package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	//	"encoding/json"
	"github.com/caarlos0/env"
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var (
	client *http.Client
	pool   *x509.CertPool
)

type envConfig struct {
	ListenPort string `env:"LISTEN_PORT" envDefault:"8080"`
}

//Config stores global env variables
var Config = envConfig{}

// Global state of test services
//var State = make(map[string]string)

var Logger = logrus.New()

func Info(args ...interface{}) {
    Logger.Info(args...)
}

var Lock = struct {
	sync.RWMutex
	State map[string]string
}{State: make(map[string]string)}

//var State = []Service{}

// func lookupService()
func lookupService(s Service) {
	response := make(chan int)
	testServiceState := "OK"

	go func() {
		resp, err := client.Head(s.URL)
		if err != nil {
			response <- http.StatusServiceUnavailable
			return
		}

		response <- resp.StatusCode
	}()

	code := 0

	select {
	case <-time.After(2 * time.Second):
		code = http.StatusServiceUnavailable
		break
	case code = <-response:
		break
	}
	if code != 200 {
		testServiceState = "Service Unavailable"
	}

	Lock.Lock()
	defer Lock.Unlock()
	Lock.State[s.Name] = testServiceState
	//	State[s.Name] = testServiceState

	//State  = Service{s.Name,s.URL}
//	Logger.Info("Probing test service ", s.Name)
}

func init() {
      Logger.Level = logrus.InfoLevel
      Logger.Formatter = &logrus.JSONFormatter{}

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}

	d := LoadDataFile("data/services.yaml")
	for _, service := range d {
		go func(s Service) {
			if s.Interval == 0 {
				s.Interval = 10
			}
			for _ = range time.Tick(time.Duration(s.Interval) * time.Second) {
				lookupService(s)
			}
		}(service)
	}
}

func main() {

	err := env.Parse(&Config)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	// Add handlers and start the server
	Address := ":" + Config.ListenPort

        gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(Logrus())
        router.GET("/services", GlobalServiceStatus)

	s := &http.Server{
		Addr:           Address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	Logger.Info("Application listening on port ", Config.ListenPort)
	s.ListenAndServe()

}
