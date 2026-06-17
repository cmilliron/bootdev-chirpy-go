package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/cmilliron/bootdev-chirpy-go/internal/auth"
	"github.com/google/uuid"
)

type UserWithToken struct{
		ID			uuid.UUID  	`json:"id"`
		CreatedAt 	time.Time	`json:"created_at"`
		UpdatedAt 	time.Time	`json:"updated_at"`
		Email		string		`json:"email"`
		Token 		string		`json:"token"`
	}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds *int `json:"expires_in_seconds,omitempty"`
	}		

	decoder := json.NewDecoder(r.Body)
	var params parameters
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}
	var expires time.Duration = 3600 * time.Second
	if params.ExpiresInSeconds != nil && (*params.ExpiresInSeconds > 0 && *params.ExpiresInSeconds < 3600)  {
			expires = time.Duration(*params.ExpiresInSeconds)
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

	token, err := auth.MakeJWT(user.ID, cfg.secret, expires)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}

	mappedUser := UserWithToken{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
		Token: token,
	}
	
	sendApiResponse(w, http.StatusOK, mappedUser)
}