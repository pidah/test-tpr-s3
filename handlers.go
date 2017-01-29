package main

import (
//	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

//func viewCountryCodes(w http.ResponseWriter, req *http.Request) {
//	rs, err := getCountryCodes()
//	if err != nil {
//		//		TODO..do not panic; use a recovery handler
//		panic(err)
//	}
//
//	bs, err := json.Marshal(rs)
//	if err != nil {
//		//		TODO..do not panic; use a recovery handler
//		panic(err)
//	}
//	w.Write(bs)
//}

type supportCORS struct {
	router *mux.Router
}

func (server *supportCORS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", `POST, GET, OPTIONS,
        	PUT, DELETE`)
		w.Header().Set("Access-Control-Allow-Headers",
			`Accept, Content-Type, Content-Length, Accept-Encoding,
            X-CSRF-Token, Authorization`)
	}
	// Stop here if its Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	server.router.ServeHTTP(w, r)
}

func servHome(w http.ResponseWriter, r *http.Request) {
	RenderTemplate(w, r, "index/home", nil)
}

//AddHandlers creates a router and adds handlers
func AddHandlers() *mux.Router {
	router := mux.NewRouter()
	http.Handle("/", &supportCORS{router})
//	router.HandleFunc("/status", viewCountryCodes).Methods("GET")
	router.HandleFunc("/", servHome).Methods("GET")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	return router
}
