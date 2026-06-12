package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func respondWithError(w http.ResponseWriter, code int, msg string, err error)  {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorRepsonse struct {
		Error string `json:"error"`
	}
    respBody := errorRepsonse{
		Error: msg + ": " + err.Error(),
    }

	sendApiResponse(w, code, respBody)
}

func respondWithSuccess(w http.ResponseWriter, code int, cleanedText string) {
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	payload := returnVals{
		CleanedBody: cleanedText,
	}

	sendApiResponse(w, code, payload)
}


func sendApiResponse(w http.ResponseWriter, code int, payload interface{}) {
    data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(data)
}