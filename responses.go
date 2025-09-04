package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type parameters struct {
	Body string `json:"body"`
}

type err_resp struct {
	Error string `json:"error"`
}

type valid_resp struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func sendValidResponse(w http.ResponseWriter, r *http.Request, cleaned_body string) {
	w.WriteHeader(200)
	respBody := valid_resp{
		Valid:       true,
		CleanedBody: cleaned_body,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	w.Write(dat)
}

func sendErrorResponse(w http.ResponseWriter, r *http.Request, err_str string) {
	w.WriteHeader(500)
	respBody := err_resp{
		Error: fmt.Sprintf("Something went wrong: %s", err_str),
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}

	w.Write(dat)
}

func sendChirpTooLong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(400)
	respBody := err_resp{
		Error: "Chirp is too long",
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}

	w.Write(dat)
}
