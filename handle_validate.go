package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)


func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body 	string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
    }

	cleanedChirp, err := validateChirp(params.Body)
	if (err != nil) {
		respondWithError(w, 400, err.Error(), err)
	} else {	
		respondWithSuccess(w, http.StatusOK, cleanedChirp)
	}
	
}

func validateChirp(chirp string) (string, error) {
	if (len(chirp) > 140) {
		return "", fmt.Errorf("Chirp is too long")
	}
	cleanedChirp, err := cleanVulgar(chirp); 
	if (err != nil) {
		return "", err
	}
	return cleanedChirp, nil
}

func cleanVulgar(chirp string) (string, error) {
	vulgarWords := map[string]struct{}{
		"kerfuffle": {}, 
		"sharbert": {}, 
		"fornax": {},
	}
	words := strings.Split(chirp," ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := vulgarWords[lowerWord]; ok {
			words[i] = "****"
		}
	
	}
	cleanedChirp := strings.Join(words, " ")
	return cleanedChirp, nil
}