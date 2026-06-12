package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cmilliron/bootdev-chirpy-go/internal/auth"
)

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
	// if result == false {
	// 	// type failure struct {
	// 	// 	Msg string `json:"msg"`
	// 	// }
	// 	// sendApiResponse(w, http.StatusUnauthorized, failure{
	// 	// 	Msg: "Incorrect email or password",
	// 	// })
	// 	respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", errors.New("Authorization error")) //change for production
	// 	return
	// }

	mappedUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
	}
	
	sendApiResponse(w, http.StatusOK, mappedUser)
}