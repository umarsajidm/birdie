package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GET ChirpsAscHandler
func (cfg *apiConfig) ChirpsAscHandler(w http.ResponseWriter, r *http.Request) {

	//fetching the chirps from db
	ChirpsAscDB, err := cfg.DB.GetChirpsAsc(r.Context())
	if err != nil {
		respondWithERROR(w, http.StatusInternalServerError, "can't get the chirpsdb in asc order")
		return
	}

	chirps := []Chirp{}

	for _, ChirpsAscDB := range ChirpsAscDB {
		chirps = append(chirps, Chirp{
			ID:        ChirpsAscDB.ID,
			CreatedAt: ChirpsAscDB.CreatedAt,
			UpdatedAt: ChirpsAscDB.UpdatedAt,
			Body:      ChirpsAscDB.Body,
			UserID:    ChirpsAscDB.UserID,
		})
	}

	respondWithValid(w, http.StatusOK, chirps)
}

// MetricsHandler
func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {

	hits := cfg.fileserverHits.Load()
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

// ResetHandler
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

// handleUsersCreate
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
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
