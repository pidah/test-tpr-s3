package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
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
var State = make(map[string]string)

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
	State[s.Name] = testServiceState

	//State  = Service{s.Name,s.URL, code}
	//fmt.Println(State)
}

func init() {

	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
	client = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{RootCAs: pool}}}

	interval := 3
	d := LoadDataFile("data/services.yaml")
	go func() {
		for _ = range time.Tick(time.Duration(interval) * time.Second) {
			for _, service := range d {
				go lookupService(service)
			}
		}
	}()
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
	serverErr := http.ListenAndServe(Address, loggedRouter)
	if serverErr != nil {
		fmt.Println(serverErr)
	}
}
