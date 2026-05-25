package main

import (
	"log"
	"net/http"
	"sync/atomic"
	"fmt"
)

type apiConfig struct{

	fileserverHits atomic.Int32
}

//middleware for statefull handler
func	(cfg *apiConfig) middlewareMatricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

//middleware Metrics
func	(cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {

		hits := cfg.fileserverHits.Load()
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

//middleware Reset
func	(cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	
		cfg.fileserverHits.Store(0)
		w.WriteHeader(http.StatusOK)
}

//middleware log
func	middlewareLog(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func	customHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") //response header
	
	w.WriteHeader(http.StatusOK) //set status code

	w.Write([]byte("OK"))

}


func main() {

	cfg := &apiConfig{}
	
	mux := http.NewServeMux() //creating a new mux
	
	fs := http.FileServer(http.Dir("."))

	stripped := http.StripPrefix("/app", fs)

	handler := cfg.middlewareMatricsInc(middlewareLog(stripped))

	mux.Handle("/app/", handler)

	mux.HandleFunc("GET /healthz/", customHandler)

	mux.HandleFunc("GET /metrics/", cfg.MetricsHandler)

	mux.HandleFunc("POST /reset/", cfg.ResetHandler)

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	//creating the new server
	server := &http.Server{
		Addr: ":80",
		Handler: mux,
	}

	log.Println("starting the server on :80")

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}