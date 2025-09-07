package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ErrResp struct {
	Error string `json:"error"`
}

type ValidResp struct {
	Valid       bool   `json:"valid"`
	CleanedBody string `json:"cleaned_body"`
}

func sendLoginAccepted(w http.ResponseWriter, user User) {
	w.WriteHeader(http.StatusOK)
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		sendErrorResponse(w, err.Error())
		return
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func sendUserNotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func sendIncorrectPasswordResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func sendChirpNotFoundResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func sendChirpResponse(w http.ResponseWriter, chirp Chirp) {
	w.WriteHeader(http.StatusOK)
	dat, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		sendErrorResponse(w, err.Error())
		return
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func sendChirpsResponse(w http.ResponseWriter, chirps []Chirp) {
	w.WriteHeader(http.StatusOK)
	dat, err := json.Marshal(chirps)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		sendErrorResponse(w, err.Error())
		return
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func sendCreatedChirpResponse(w http.ResponseWriter, chirp Chirp) {
	w.WriteHeader(http.StatusCreated)
	dat, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		sendErrorResponse(w, err.Error())
		return
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func sendCleanedResponse(w http.ResponseWriter, cleaned_body string) {
	w.WriteHeader(http.StatusOK)
	respBody := ValidResp{
		Valid:       true,
		CleanedBody: cleaned_body,
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, err_str string) {
	w.WriteHeader(http.StatusInternalServerError)
	respBody := ErrResp{
		Error: fmt.Sprintf("Something went wrong: %s", err_str),
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}

func sendChirpTooLong(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	respBody := ErrResp{
		Error: "Chirp is too long",
	}

	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}

	_, err = w.Write(dat)

	if err != nil {
		log.Println("Error writing response: %w", err)
	}
}
