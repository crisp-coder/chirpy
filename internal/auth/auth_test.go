package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT_Correct(t *testing.T) {
	userId, _ := uuid.Parse("3f1c2e7a-9b5b-4c2a-8f6a-1a2b3c4d5e6f")
	secret := "test-secret"
	token, err := MakeJWT(userId, secret, time.Minute)
	if err != nil {
		t.Error("token creation failed")
	}

	gotId, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("validation failed, err not nil")
	}

	if gotId != userId {
		t.Fatalf("validation failed, wrong id in token, %v", gotId)
	}
}

func TestJWT_Expired(t *testing.T) {
	userId, _ := uuid.Parse("3f1c2e7a-9b5b-4c2a-8f6a-1a2b3c4d5e6f")
	secret := "test-secret"
	token, err := MakeJWT(userId, secret, -1*time.Second)
	if err != nil {
		t.Fatalf("token creation failed")
	}

	gotId, err := ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("validation failed, err is nil for expired token")
	}

	if gotId != uuid.Nil {
		t.Fatalf("validation failed, wrong id in token, %v", gotId)
	}
}

func TestGetBearerToken(t *testing.T) {
	userId, _ := uuid.Parse("3f1c2e7a-9b5b-4c2a-8f6a-1a2b3c4d5e6f")
	secret := "test-secret"
	token, err := MakeJWT(userId, secret, time.Minute)
	if err != nil {
		t.Error("token creation failed")
	}
	req, _ := http.NewRequest(http.MethodGet, "https://localhost:8080/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	bearerToken, err := GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("error parsing bearer token")
	}

	if bearerToken != token {
		t.Fatalf("bearer token does not match jwt token")
	}
}
