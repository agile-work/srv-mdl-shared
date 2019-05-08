package shared

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	addr = flag.String("port", "", "TCP port to listen to")
	cert = flag.String("cert", "cert.pem", "Path to certification")
	key  = flag.String("key", "key.pem", "Path to certification key")
)

//ListenAndServe default module api listen and server
func ListenAndServe(port string, moduleRouter *chi.Mux) {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	flag.Parse()
	if *addr != "" {
		port = *addr
	}
	if port == "" {
		panic("Invalid module port")
	}

	caCert, err := ioutil.ReadFile(*cert)
	if err != nil {
		panic("Invalid service certificate")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	router := chi.NewRouter()
	router.Use(
		middleware.Heartbeat("/ping"),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
	)
	router.Mount("/api/v1", moduleRouter)

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		TLSConfig:    tlsConfig,
	}

	go func() {
		fmt.Printf("Service listening on %s\n", port)
		if err := srv.ListenAndServeTLS(*cert, *key); err != nil {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan
	fmt.Println("Shutting down Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	fmt.Println("Service stopped!")
}
