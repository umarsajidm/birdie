package main

import (
	"log"
	"net/http"
)

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

	mux := http.NewServeMux() //creating a new mux
	
	
	fs := http.FileServer(http.Dir("."))

	stripped := http.StripPrefix("/app", fs)

	mux.Handle("/app/", middlewareLog(stripped))

	mux.HandleFunc("/healthz", customHandler)

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