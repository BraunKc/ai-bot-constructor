package userdomain

import (
	"errors"
	"strings"
)

var (
	ErrEmptyUsername         = errors.New("empty username")
	ErrUsernameMustBeLonger  = errors.New("username must be longer")
	ErrUsernameMustBeShorter = errors.New("username must be shorter")

	ErrEmptyPassword        = errors.New("empty password")
	ErrPasswordMustBeLonger = errors.New("password must be longer")
)

type Hasher interface {
	Hash(str string) (string, error)
	Compare(hash, str string) bool
}

type Username string

func NewUsername(username string) (Username, error) {
	trimmed := strings.TrimSpace(username)
	lenTrimmed := len(trimmed)

	if lenTrimmed < 3 {
		return "", ErrUsernameMustBeLonger
	}
	if lenTrimmed > 32 {
		return "", ErrUsernameMustBeShorter
	}

	return Username(username), nil
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

func (ph PasswordHash) String() string {
	return string(ph)
}
