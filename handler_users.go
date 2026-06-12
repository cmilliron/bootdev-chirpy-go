package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email 	string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	var data parameters
	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	user, err :=  cfg.db.CreateUser(r.Context(), data.Email)
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