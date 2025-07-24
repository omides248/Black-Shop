package identity

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidPassword    = errors.New("invalid password")
)

type UserID string

type User struct {
	ID           UserID
	Name         string
	Email        string
	PasswordHash string
}

func NewUser(name, email, plainPassword string) (*User, error) {

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashPassword),
	}, nil
}

func (u *User) CheckPassword(plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidPassword
		}
		return err
	}
	return nil
}
