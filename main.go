package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gorilla/handlers"
	//        "gopkg.in/gin-gonic/gin.v1"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"
)

type envConfig struct {
	ListenPort string `env:"LISTEN_PORT" envDefault:"8080"`
	ConnectURL string `env:"CONNECT_URL" envDefault:"127.0.0.1"`
}

//Config stores global env variables
var Config = envConfig{}

// func lookupService()
func lookupService(url string, errChan chan error) {
	response := make(chan struct{})

	go func() {
		resp, err := http.Get(url)

		if err != nil {
			errChan <- err
		}

		// TODO: check resp
		// ...

		response <- struct{}{}
	}()

	select {
	case <-time.After(2 * time.Second):
		errChan <- errors.New("timeout for service ")
		return
	case <-response:
		return
	}
}

// =============
// HTTP endpoint
//return gin.JSON(Status) // global Status

func init() {
	d := LoadDataFile("services.json")
	go func() {
		for _ = range time.Tick(10 * time.Second) {
			var wg sync.WaitGroup
			var Status error

			// N responses
			wg.Add(len(d))
			errChan := make(chan error, len(d)) // buffered for N responses

			for _, service := range d {
				go func() {
					lookupService(service.(string), errChan)
					wg.Done()
				}()
			}

			wg.Wait()

			// Here all of the goroutines somehow terminated
			for err := range errChan {
				Status = err
				break
			}

			Status = nil
		}
	}()
}

func main() {

	err := env.Parse(&Config)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	//	d := LoadDataFile("services.json")

	//        fmt.Println("services file successfully loaded")

	// Add handlers and start the server
	Address := ":" + Config.ListenPort

	loggedRouter := handlers.LoggingHandler(os.Stdout, AddHandlers())

	fmt.Println("Application listening on port", Config.ListenPort)
	serverErr := http.ListenAndServe(Address, loggedRouter)
	if serverErr != nil {
		fmt.Println(serverErr)
	}
}
