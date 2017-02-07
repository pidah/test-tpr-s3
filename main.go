package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	//	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/caarlos0/env"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
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
	//fmt.Println(State)
}

func init() {

	logger := logrus.New()
	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.JSONFormatter{}

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}

	d := LoadDataFile("data/services.yaml")
	for _, service := range d {
		go func(s Service) {
			if s.Internal == 0 {
				s.Internal = 10
			}
			for _ = range time.Tick(time.Duration(s.Interval) * time.Second) {
				fmt.Println(s.Interval)
				lookupService(service)
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

	loggedRouter := handlers.LoggingHandler(os.Stdout, AddHandlers())

	fmt.Println("Application listening on port", Config.ListenPort)
	serverErr := http.ListenAndServe(Address, nil)
	if serverErr != nil {
		fmt.Println(serverErr)
	}

	//	http.Handle("/", L.Handler(http.HandlerFunc(exhandler), "homepage"))

	http.ListenAndServe(":8080", loggedRouter)

}
