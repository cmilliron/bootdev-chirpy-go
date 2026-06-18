package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cmilliron/bootdev-chirpy-go/internal/auth"
	"github.com/cmilliron/bootdev-chirpy-go/internal/database"
	"github.com/google/uuid"
)

type UserWithToken struct{
		ID				uuid.UUID  	`json:"id"`
		CreatedAt 		time.Time	`json:"created_at"`
		UpdatedAt 		time.Time	`json:"updated_at"`
		Email			string		`json:"email"`
		Token 			string		`json:"token"`
		RefreshToken 	string		`json:"refresh_token"`
	}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	string `json:"email"`
		Password string `json:"password"`
	}		

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	fmt.Println(user)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find user", err) //change for production
		return
	}

	result, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !result {
		if !result { 
			err = errors.New("Authorization error")
		}
		respondWithError(w, http.StatusUnauthorized, "Something went wrong with your login", err) //change for production
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, cfg.defaults.tokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	newRefreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().UTC().Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	mappedUser := UserWithToken{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: newRefreshToken.Token,
	}
	
	sendApiResponse(w, http.StatusOK, mappedUser)
}


func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No Token", err)
		return
    }
	
	tokenData, err := cfg.db.GetRefreshToken(r.Context(), refresh_token)
	if err != nil ||  tokenData.RevokedAt.Valid || !tokenData.ExpiresAt.After(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token expired or revoked", err)
		return
	}
	
	
	accessToken, err := auth.MakeJWT(tokenData.UserID, cfg.secret, cfg.defaults.tokenExpiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}
	
	type response struct {
		Token	string `json:"token"`
	}
	
	sendApiResponse(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {

		respondWithError(w, http.StatusUnauthorized, "No Token", err)
		return
    }

	_, err = cfg.db.UpdateRevokeRefreshToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	sendApiResponse(w, http.StatusNoContent, nil)
	
}  

