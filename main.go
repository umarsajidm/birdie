package main

import (
	"database/sql"
	// "encoding/json"
	"os"
	"log"
	"net/http"
	"sync/atomic"
	// "fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/umarsajidm/birdie/internal/database"
)

type User struct {
	ID	uuid.UUID `json:"id"`
	CreatedAt time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	Email	string	`json:"email"`
}


type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
	Platform		string
}

// middleware for statefull handler
func (cfg *apiConfig) middlewareMatricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// middleware log
func middlewareLog(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func customHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8") //response header

	w.WriteHeader(http.StatusOK) //set status code

	w.Write([]byte("OK"))

}

func main() {

	//load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//grab the connection string
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	//open the connection to postgres
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	//initialize your sqlc queries and add them to config
	dbQueries := database.New(db)

	cfg := &apiConfig{
		DB: dbQueries,
		Platform: os.Getenv("PLATFORM"),
	}

	mux := http.NewServeMux() //creating a new mux

	fs := http.FileServer(http.Dir("."))

	stripped := http.StripPrefix("/app", fs)

	handler := cfg.middlewareMatricsInc(middlewareLog(stripped))

	mux.Handle("/app/", handler)

	mux.HandleFunc("GET /healthz", customHandler)

	mux.HandleFunc("GET /metrics", cfg.MetricsHandler)
	
	mux.HandleFunc("GET /api/chirps", cfg.ChirpsAscHandler)

	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.ChirpByIDHandler)

	mux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	mux.HandleFunc("POST /api/chirps", cfg.handleChirpsCreate)

	mux.HandleFunc("POST /api/users", cfg.handleUsersCreate)


	//creating the new server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("starting the server on :8080")

	err = server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
