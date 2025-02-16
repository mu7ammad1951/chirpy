package auth

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret")

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {

	cases := []struct {
		name          string
		headerValue   string
		expected      string
		wantErr       bool
		missingHeader bool
	}{
		{
			name:        "Valid Header",
			headerValue: "Bearer TOKEN_STRING",
			expected:    "TOKEN_STRING",
			wantErr:     false,
		},
		{
			name:        "Missing Token String",
			headerValue: "Bearer",
			expected:    "",
			wantErr:     true,
		},
		{
			name:        "Malformed Header Value",
			headerValue: "BearerTOKEN_STRING",
			expected:    "",
			wantErr:     true,
		},
		{
			name:        "Missing Header Value",
			headerValue: "",
			expected:    "",
			wantErr:     true,
		},
		{
			name:          "Missing Header",
			headerValue:   "",
			expected:      "",
			wantErr:       true,
			missingHeader: true,
		},
		{
			name:        "Malformed Header Value II",
			headerValue: "Bearer         ",
			expected:    "",
			wantErr:     true,
		},
	}

	for _, tc := range cases {
		dummyReq, err := http.NewRequest("POST", "http://localhost:8080/api/chirps", nil)
		if err != nil {
			t.Fatalf("error creating request: %v", err)
		}
		if !tc.missingHeader {
			dummyReq.Header.Add("Authorization", tc.headerValue)
		}

		expect, err := GetBearerToken(dummyReq.Header)
		if (err != nil) != tc.wantErr || expect != tc.expected {
			t.Errorf("test: '%v' failed - expected: %v, got: %v, wantErr: %v, err: %v", tc.name, tc.expected, expect, tc.wantErr, err)
		}
	}

}
