package main

import (
	"encoding/json"
	"net/http"
	"log"
)

func respondWithValid(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func respondWithERROR(w http.ResponseWriter, code int, s string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithValid(w, code, errorResponse{Error: s})
}

func	validateChirpHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithERROR(w, http.StatusBadRequest, "something went wrong in decoding")
		return
	}

	if len(params.Body) > 140 {
		respondWithERROR(w, http.StatusBadRequest, "chirp is too long")
		return 
	}

	type validResponse struct {
		Valid bool "json:valid"
	}

	respondWithValid(w, http.StatusOK, validResponse{Valid: true})

}