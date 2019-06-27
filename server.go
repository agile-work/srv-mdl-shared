package shared

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/agile-work/srv-mdl-shared/middlewares"

	"github.com/agile-work/srv-mdl-shared/models/translation"
	"github.com/agile-work/srv-shared/rdb"
	"github.com/agile-work/srv-shared/service"
	"github.com/agile-work/srv-shared/socket"
	"github.com/agile-work/srv-shared/util"

	"github.com/agile-work/srv-shared/constants"

	"gopkg.in/go-playground/validator.v9"

	"github.com/agile-work/srv-shared/sql-builder/db"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var (
	cert       = flag.String("cert", "cert.pem", "Path to certification")
	key        = flag.String("key", "key.pem", "Path to certification key")
	mdlHost    = flag.String("host", "", "Module host")
	mdlPort    = flag.Int("port", -1, "Module port")
	dbHost     = flag.String("dbHost", "cryo.cdnm8viilrat.us-east-2.rds-preview.amazonaws.com", "Database host")
	dbPort     = flag.Int("dbPort", 5432, "Database port")
	dbUser     = flag.String("dbUser", "cryoadmin", "Database user")
	dbPassword = flag.String("dbPassword", "x3FhcrWDxnxCq9p", "Database password")
	dbName     = flag.String("dbName", "cryo", "Database name")
	redisHost  = flag.String("redisHost", "localhost", "Redis host")
	redisPort  = flag.Int("redisPort", 6379, "Redis port")
	redisPass  = flag.String("redisPass", "redis123", "Redis password")
	wsHost     = flag.String("wsHost", "localhost", "Realtime host")
	wsPort     = flag.Int("wsPort", 8010, "Realtime port")
)

// Validate global instance of the validator
var Validate *validator.Validate

// ListenAndServe default module api listen and server
func ListenAndServe(name, host string, port int, moduleRouter *chi.Mux) {
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	flag.Parse()
	if *mdlPort != -1 {
		port = *mdlPort
	}
	if *mdlHost != "" {
		host = *mdlHost
	}

	pid := os.Getpid()
	module := service.New(name, constants.ServiceTypeModule, host, port, pid)

	fmt.Printf("Starting Module %s...\n", module.Name)
	fmt.Printf("[Instance: %s | PID: %d]\n", module.InstanceCode, module.PID)

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
	fmt.Println("Database connected")

	rdb.Init(*redisHost, *redisPort, *redisPass)
	defer rdb.Close()

	socket.Init(module, *wsHost, *wsPort)
	defer socket.Close()

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
		Addr:         module.URL(),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
		TLSConfig:    tlsConfig,
	}

	Validate = validator.New()

	go func() {
		fmt.Printf("Service %s listening on %d\n", name, port)
		if err := httpServer.ListenAndServeTLS(*cert, *key); err != nil {
			if strings.Contains(err.Error(), "bind: address already in use") {
				if err := rdb.Delete("module:def:" + module.InstanceCode); err != nil {
					fmt.Println(err)
				}
				if _, err := rdb.LRem("api:modules", 0, module.InstanceCode); err != nil {
					fmt.Println(err)
				}
				fmt.Println("")
				log.Fatalf("port %d already in use\n", port)
			}
		}
	}()

	rdb.LPush("api:modules", module.InstanceCode)
	rdb.Set("module:def:"+module.InstanceCode, module.JSON(), 0)

	deadline := time.Now().Add(15 * time.Second)
	for {
		if socket.Available() || time.Now().After(deadline) {
			if err := socket.Emit(socket.Message{
				Recipients: []string{"service.api"},
				Data:       "reload",
			}); err != nil {
				fmt.Println(err)
			}
			break
		}
	}

	<-stopChan
	fmt.Println("\nShutting down Service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	httpServer.Shutdown(ctx)
	if err := rdb.Delete("module:def:" + module.InstanceCode); err != nil {
		fmt.Println(err)
	}
	if _, err := rdb.LRem("api:modules", 0, module.InstanceCode); err != nil {
		fmt.Println(err)
	}
	if err := socket.Emit(socket.Message{
		Recipients: []string{"service.api"},
		Data:       "reload",
	}); err != nil {
		fmt.Println(err)
	}
	defer cancel()
	fmt.Println("Service stopped!")
}
