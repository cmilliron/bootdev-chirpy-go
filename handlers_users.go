package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cmilliron/bootdev-chirpy-go/internal/auth"
	"github.com/cmilliron/bootdev-chirpy-go/internal/database"
	"github.com/google/uuid"
)

// handlerCreateUser registers a new user account with a hashed password.
func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	var data parameters
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          data.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	mappedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.CreatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	sendApiResponse(w, http.StatusCreated, mappedUser)

}

// handlerUpdateUser updates the authenticated user's email and password.
func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var jsonRequest parameters
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&jsonRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No access token", err)
		return
	}

	userId, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token malformed or ", err)
		return
	}

	hashedPassword, err := auth.HashPassword(jsonRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          jsonRequest.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't update user", err)
		return
	}

	mappedUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.CreatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	sendApiResponse(w, http.StatusOK, mappedUser)
}

// handlerUpdateChirpyRed processes the Polka webhook event and marks a user
// as Chirpy Red when their upgrade event is received.
func (cfg *apiConfig) handlerUpdateChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil || apiKey != cfg.polkaKey {
		output := fmt.Sprintf("User not authorized: %v", err)
		respondWithError(w, http.StatusUnauthorized, output, err)
		return
	}

	var params parameters
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		sendApiResponse(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.UpdateUserPremiumStatus(r.Context(), database.UpdateUserPremiumStatusParams{
		IsChirpyRed: true,
		ID:          params.Data.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User Not Found", err)
		return
	}

	sendApiResponse(w, http.StatusNoContent, nil)
}
