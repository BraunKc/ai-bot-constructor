package userdomain

import "github.com/google/uuid"

type User struct {
	id           uuid.UUID
	username     Username
	passwordHash PasswordHash
}

func NewUser(username, password string, hasher Hasher) (*User, error) {
	u, err := NewUsername(username)
	if err != nil {
		return nil, err
	}

	ph, err := NewPasswordHash(password, hasher)
	if err != nil {
		return nil, err
	}

	return &User{
		id:           uuid.New(),
		username:     u,
		passwordHash: ph,
	}, nil
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) Username() Username {
	return u.username
}

func (u *User) CheckPassword(password string, hasher Hasher) bool {
	return u.passwordHash.Compare(password, hasher)
}

func (u *User) UpdateUsername(newUsername string) error {
	if u.username.String() == newUsername {
		return nil
	}

	nu, err := NewUsername(newUsername)
	if err != nil {
		return err
	}

	u.username = nu

	return nil
}
