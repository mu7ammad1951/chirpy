package auth

import "testing"

func TestHashCompareAndPasswordNormalValid(t *testing.T) {
	password := "helloworldhelloworld"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`failed to hash password: %v`, err)
	}
	err = CheckPasswordHash(password, passwordHash)
	if err != nil {
		t.Fatalf(`password hash check failed: %v`, err)
	}
}

func TestHashCompareAndPasswordNormalInvalid(t *testing.T) {
	password := "helloworldhelloworld"
	passwordHash, err := HashPassword(password)
	if err != nil {
		t.Fatalf(`failed to hash password: %v`, err)
	}
	err = CheckPasswordHash("worldhellohelloworld", passwordHash)
	if err == nil {
		t.Fatal(`password hash check failed: allowed invalid password`)
	}
}

func TestHashCompareAndPasswordBlank(t *testing.T) {
	password := ""
	_, err := HashPassword(password)
	if err == nil {
		t.Fatal(`allowed hashing blank password`)
	}
}
