package main

import (
	"encoding/json"
	"net/http"

	"github.com/cmilliron/bootdev-chirpy-go/internal/auth"
	"github.com/cmilliron/bootdev-chirpy-go/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	string `json:"email"`
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

	user, err :=  cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: data.Email, 
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	mappedUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
	}
	sendApiResponse(w, http.StatusCreated, mappedUser)

}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 		string `json:"email"`
		Password	string `json:"password"`
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
		Email: jsonRequest.Email,
		HashedPassword: hashedPassword,
		ID: userId,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't update user", err)
		return
    }

	mappedUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
	}
	sendApiResponse(w, http.StatusOK, mappedUser)

	

}