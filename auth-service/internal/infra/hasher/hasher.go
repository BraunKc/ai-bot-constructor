package hasherinfra

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	Cost int
}

func (h *Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hash from string: %w", err)
	}

	return string(hash), nil
}

func (h *Hasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
