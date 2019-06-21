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

	"github.com/agile-work/srv-mdl-shared/middlewares"

	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/util"

	"github.com/agile-work/srv-shared/constants"

	"github.com/agile-work/srv-shared/service"

	"gopkg.in/go-playground/validator.v9"

	"github.com/agile-work/srv-shared/sql-builder/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	moduleName = flag.String("name", "", "Name of this module instance")
	addr       = flag.String("port", "", "TCP port to listen to")
	cert       = flag.String("cert", "cert.pem", "Path to certification")
	key        = flag.String("key", "key.pem", "Path to certification key")
	dbHost     = flag.String("dbHost", "cryo.cdnm8viilrat.us-east-2.rds-preview.amazonaws.com", "Database host")
	dbPort     = flag.Int("dbPort", 5432, "Database port")
	dbUser     = flag.String("dbUser", "cryoadmin", "Database user")
	dbPassword = flag.String("dbPassword", "x3FhcrWDxnxCq9p", "Database password")
	dbName     = flag.String("dbName", "cryo", "Database name")
)

// Validate global instance of the validator
var Validate *validator.Validate

// ListenAndServe default module api listen and server
func ListenAndServe(name, port string, moduleRouter *chi.Mux) {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	flag.Parse()
	if *addr != "" {
		port = *addr
	}
	if port == "" {
		panic("Invalid module port")
	}

	if *moduleName != "" {
		name = *moduleName
	}
	if name == "" {
		name = "undefined"
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

	err = db.Connect(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName, false)
	if err != nil {
		panic("Database error")
	}
	defer db.Close()

	ws, err := service.Register(name, constants.ServiceTypeModule)
	if err != nil {
		fmt.Printf("\nUnable to connect to realtime socket. Error: %s", err.Error())
	}

	params, err := util.GetSystemParams()
	if err == nil {
		translation.SystemDefaultLanguageCode = params[constants.SysParamLanguageCode]
	}

	router := chi.NewRouter()
	router.Use(
		middleware.Heartbeat("/ping"),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middlewares.Translation,
	)
	router.Mount("/api/v1", moduleRouter)

	httpServer := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		TLSConfig:    tlsConfig,
	}

	Validate = validator.New()

	go func() {
		fmt.Printf("Service listening on %s\n", port)
		if err := httpServer.ListenAndServeTLS(*cert, *key); err != nil {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan
	fmt.Println("\nShutting down Service...")
	ws.Down()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	httpServer.Shutdown(ctx)
	defer cancel()
	fmt.Println("Service stopped!")
}
