package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/gorilla/handlers"
	//        "gopkg.in/gin-gonic/gin.v1"
	"net/http"
	"os"
	"time"
)

type envConfig struct {
	ListenPort string `env:"LISTEN_PORT" envDefault:"8080"`
}

//Config stores global env variables
var Config = envConfig{}

// Global state of test services
var State = make(map[string]int)
//var State = []Service{}

// func lookupService()
func lookupService(s Service) {
	response := make(chan int)

	go func() {
		resp, err := http.Head(s.URL)
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
        State[s.Name] = code
//State  = Service{s.Name,s.URL, code}
	fmt.Println(State)
}

// =============
// HTTP endpoint
//return gin.JSON(Status) // global Status

func init() {
	d := LoadDataFile("services.json")
	go func() {
		for _ = range time.Tick(3  * time.Second) {
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
