package main

import (
	"encoding/json"
	//	"fmt"
	"github.com/sirupsen/logrus"
        "gopkg.in/gin-gonic/gin.v1"
	"net/http"
        "time"
)

func GlobalServiceStatus(c *gin.Context) {
	Lock.RLock()
	defer Lock.RUnlock()

	for _, value := range Lock.State {
		if value != "OK" {
			c.Writer.WriteHeader(http.StatusServiceUnavailable)
		}
	}
	bs, err := json.Marshal(Lock.State)
	if err != nil {
		//		TODO..do not panic; use a recovery handler
		panic(err)
	}
	c.Writer.Write(bs)
}

//     func SingleServiceStatus(c *gin.Context) {
//	Lock.RLock()
//	defer Lock.RUnlock()
//
//	vars := mux.Vars(req)
//	ts := vars["testService"]
//	if val, ok := Lock.State[ts]; !ok {
//		w.WriteHeader(http.StatusNotFound)
//		w.Write([]byte("☄ hey!, the requested test service could not be found."))
//
//	} else {
//		bs, err := json.Marshal(val)
//		if err != nil {
//			//		TODO..do not panic; use a recovery handler
//			panic(err)
//		}
//		if val != "OK" {
//			w.WriteHeader(http.StatusServiceUnavailable)
//		}
//		w.Write(bs)
//	}
//}

//type supportCORS struct {
//	router *mux.Router
///}

//func (server *supportCORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if origin := r.Header.Get("Origin"); origin != "" {
//		w.Header().Set("Access-Control-Allow-Origin", origin)
//		w.Header().Set("Access-Control-Allow-Methods", `POST, GET, OPTIONS,
  //      	PUT, DELETE`)
//		w.Header().Set("Access-Control-Allow-Headers",
//			`Accept, Content-Type, Content-Length, Accept-Encoding,
  //          X-CSRF-Token, Authorization`)
//	}
//	// Stop here if its Preflighted OPTIONS request
//	if r.Method == "OPTIONS" {
//		return
//	}
//	// Lets Gorilla work
//	server.router.ServeHTTP(w, r)
//}

func servHome(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/home", nil)
}

func Logrus() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.Request.URL.Path
		c.Next()
		end := time.Now().UTC()
		latency := end.Sub(start)
		Logger.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"ip":         c.ClientIP(),
			"duration":   latency,
			"user_agent": c.Request.UserAgent(),
		}).Info()
	}
}




//	router := mux.NewRouter()
//	http.Handle("/", &supportCORS{router})
//	router.HandleFunc("/services", (http.HandlerFunc(G), "homepage")).Methods("GET")
//	router.HandleFunc("/service/{testService}", SingleServiceStatus).Methods("GET")
//	router.HandleFunc("/", servHome).Methods("GET")
//	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
//	return router
//}
