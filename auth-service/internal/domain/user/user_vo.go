package userdomain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyUsername        = errors.New("empty username")
	ErrUsernameMustBeLonger = errors.New("username must be longer")

	ErrEmptyPassword        = errors.New("empty password")
	ErrPasswordMustBeLonger = errors.New("password must be longer")
)

type Hasher interface {
	Hash(str string) (string, error)
	Compare(hash, str string) bool
}

type Username string

func NewUsername(username string) (Username, error) {
	u := Username(username)
	if err := u.validate(); err != nil {
		return "", err
	}

	return u, nil
}

func (u Username) validate() error {
	trimmed := strings.TrimSpace(u.String())

	if trimmed == "" {
		return ErrEmptyUsername
	}

	if len(trimmed) < 3 {
		return ErrUsernameMustBeLonger
	}

	return nil
}

func (u Username) String() string {
	return string(u)
}

type PasswordHash string

func NewPasswordHash(password string, hasher Hasher) (PasswordHash, error) {
	if password == "" {
		return "", ErrEmptyPassword
	}

	if len(password) < 6 {
		return "", ErrPasswordMustBeLonger
	}

	hash, err := hasher.Hash(password)
	if err != nil {
		return "", err
	}

	return PasswordHash(hash), nil
}

func (ph PasswordHash) Compare(password string, hasher Hasher) bool {
	return hasher.Compare(string(ph), password)
}
