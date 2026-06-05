package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"log"
	"net/http"
	"sync/atomic"
	"fmt"
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

// middleware Metrics
func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {

	hits := cfg.fileserverHits.Load()
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

//ResetHandler
func (cfg *apiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {

	if cfg.Platform != "dev" {
		respondWithERROR(w, http.StatusForbidden, "reset is only allowed in dev environment")
		return
	}

	//wiping the fileserver hits
	cfg.fileserverHits.Store(0)
	//executing the sqlc delete query
	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithERROR(w, http.StatusInternalServerError, "couldnt delete users")
		return
	}
	respondWithValid(w, http.StatusOK, struct {

		Message string `json:"message"`
	}{
		Message: "hits reset and database wiped cleanly",
	})
}

//handleUsersCreate
func (cfg *apiConfig) handleUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}
	//decoding the incoming json body/ request
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithERROR(w, http.StatusBadRequest, "couldn't decode the params for creating users")
		return
	}
	//executing the sqlc query
	//passing the r.Context() to handle the request timeouts natively
	user, err := cfg.DB.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithERROR(w, http.StatusBadRequest, "couldnt create user")
		return
	}
	//map the database.User to my custom User struct and send the json response
	respondWithValid(w, http.StatusCreated, User{
		ID:	user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:	user.Email,
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

	mux.HandleFunc("POST /admin/reset", cfg.ResetHandler)

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

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
