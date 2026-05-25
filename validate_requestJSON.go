package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func	getCleanedBody(body string) string {
	words := strings.Split(body, " ")
	
	badWords := map[string]bool{
		"kerfuffle": true,
		"sharbert": true,
		"fornax": true,
	}

	for i, word := range words {
		lowercase := strings.ToLower(word)
		uppercase := strings.ToUpper(word)
		if badWords[lowercase] || badWords[uppercase] {
			
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func respondWithValid(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json %v", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
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

	cleaned := getCleanedBody(params.Body)

	type validResponse struct {
		CleanBody string `json:"cleanBody"`
	}

	respondWithValid(w, http.StatusOK, validResponse{CleanBody: cleaned})

}