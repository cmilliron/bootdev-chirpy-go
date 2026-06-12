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

	HashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}

	user, err :=  cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: data.Email, 
		HashedPassword: HashedPassword,
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