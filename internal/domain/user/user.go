package user

import (
	"errors"
	"uuid"
)

var ErrEmptyName = errors.New("client name is empty")

type User struct {
	id   uuid.UUID
	name string
}

func CreateUser(id uuid.UUID, name string) (*User, error) {
	if name == "" {
		return &User{}, ErrEmptyName
	}
	return &User{id: id, name: name}, nil
}

func RestoreUser(id uuid.UUID, name string) *User {
	return &User{id: id, name: name}
}
