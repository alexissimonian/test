package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidJwtCreationAndValidation(t *testing.T) {
	userID := uuid.New()
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	token, err := MakeJWT(userID, "secret", duration)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	resultUUID, err := ValidateJWT(token, "secret")
	if resultUUID != userID {
		t.Errorf("Expected %v, got %v\n", userID, resultUUID)
	}
}

func TestInvalidJWTSecret(t *testing.T) {
	userID := uuid.New()
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	token, err := MakeJWT(userID, "secret", duration)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	resultUUID, err := ValidateJWT(token, "anothersecret")
	if err == nil {
		t.Errorf("Expected error, got %v\n", resultUUID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	duration, err := time.ParseDuration("1ms")
	if err != nil {
		t.Errorf("%v\n", err)
	}

	token, err := MakeJWT(userID, "secret", duration)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	resultUUID, err := ValidateJWT(token, "secret")
	if err == nil {
		t.Errorf("Expected error, got %v\n", resultUUID)
	}
}

func TestGetTokenBearer(t *testing.T){
	header := http.Header{}
	header.Add("Authorization", "Bearer token")
	token, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("%v\n", err)
	}

	if token != "token" {
		t.Errorf("Expected %v, got %v\n", "token", token)
	}
}

func TestGetTokenBearerWithoutAuthorizationHeader(t *testing.T){
	header := http.Header{}
	token, err := GetBearerToken(header)
	if err == nil {
		t.Errorf("Expected error. Got %v\n", token)
	}
}