package main

import (
	"log"
	"net/http"
)

// type Server struct {
// 	Addr	string
// 	Handler	handler
// 	DisableGeneralOptionsHandler	bool
// 	TLSConfig	*tls.Config
// 	ReadTimeout	time.Duration
// 	ReadHeaderTimeout	time.Duration
// 	WriteTimeout	time.Duration
// 	IdleTimeout	time.Duration
// 	MaxHeaderBytes	int
// 	TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
// 	ConnState func(net.Conn, ConnState)
// 	ErrorLog	*log.logger
// 	BaseContext	func(net.Listener) context.Context
// 	ConnContext	func(ctx context.Context, c net.Conn) context.Context
// 	HTTP2	*HTTP2Config
// 	Protocols	*Protocols
// }


func	customHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") //response header
	
	w.WriteHeader(http.StatusOK) //set status code

	w.Write([]byte("OK"))

}


func main() {

	mux := http.NewServeMux() //creating a new mux
	
	
	fs := http.FileServer(http.Dir("."))

	stripped := http.StripPrefix("/app", fs)

	mux.Handle("/app/", stripped)

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