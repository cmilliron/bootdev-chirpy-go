package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	err = validateChirp(params.Body)
	if (err != nil) {
		respondWithError(w, 400, err.Error(), err)
	} else {	
		respondWithSuccess(w, http.StatusOK, true)
	}
	
}

func validateChirp(chirp string) error {
	if (len(chirp) > 140) {
		return fmt.Errorf("Chirp is too long")
	}
	return nil
}