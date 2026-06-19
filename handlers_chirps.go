package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/cmilliron/bootdev-chirpy-go/internal/auth"
	"github.com/cmilliron/bootdev-chirpy-go/internal/database"
	"github.com/google/uuid"
)

// POST /api/chirps
func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body 	string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Token", err)
		return
    }

	userId, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token", err)
		return
    }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
		UserID: userId,
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

// GET /api/chirps
// optional query pararameters: author_id
func (cfg *apiConfig) handlerGetAllChrips(w http.ResponseWriter, r *http.Request) {
	authorId, err := getAuthorFromQueryParams(r)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
		return
	}
	
	var chirps []database.Chirp

	if authorId == uuid.Nil {
		chirps, err = cfg.db.GetAllChirps(r.Context())
	} else {
		chirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorId)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching chirps", err)
		return
	}

	resChirps := []ChirpResponse{}
	for _, chirp := range chirps {
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
	chirpId := r.PathValue("chirpId")
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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No access token", err)
		return
    }
	
	userId, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token malformed or ...", err)
		return
    }
	
	chirpId := r.PathValue("chirpId")
	chirpUuid, err := uuid.Parse(chirpId)
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpUuid)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could fetch tweet", err)
		return
    }
	
	if chirp.UserID != userId {
		respondWithError(w, http.StatusForbidden, "Chirp does not belong to user", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "You can't delete this chirp", err)
		return
    }

	sendApiResponse(w, http.StatusNoContent, nil)
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

func getAuthorFromQueryParams(r *http.Request) (uuid.UUID, error) {
	authorIdString := r.URL.Query().Get("author_id")
	if authorIdString == "" {
		return uuid.Nil, nil
	} 
	parsedId, err := uuid.Parse(authorIdString)
	if err != nil {
		return uuid.Nil, err	
	}
	return parsedId, nil
}

