package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cmilliron/bootdev-chirpy-go/internal/database"
	"github.com/google/uuid"
)

// POST /api/chirps
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body 	string `json:"body"`
		UserId	uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
    }

	validChirp, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error validating Chirp", err)
		return
    }
	
	newChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: validChirp,
		UserID: params.UserId,
	})
		
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
    }
	
	resChirp := ChirpResponse{
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: newChirp.Body,
		UserID: newChirp.UserID,
	}
	sendApiResponse(w, http.StatusCreated, resChirp)
}

func (cfg *apiConfig) handlerGetAllChrips(w http.ResponseWriter, r *http.Request) {
	chrips, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudld fetch tweets", err)
		return
    }

	resChirps := []ChirpResponse{}
	for _, chirp := range chrips {
		resChirps = append(resChirps, ChirpResponse{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		})	
	}
	sendApiResponse(w, http.StatusOK, resChirps)
}

func (cfg *apiConfig) handlerSingleChirp(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("ChirpId")
	chirpUuid, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
    }

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Coudld fetch tweet", err)
		return
    }

	resChirp := ChirpResponse{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	sendApiResponse(w, http.StatusOK, resChirp)
}


// Helper Functions
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

