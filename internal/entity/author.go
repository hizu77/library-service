package entity

import "errors"

type Author struct {
	ID   string
	Name string
}

var (
	ErrAuthorNotFound      = errors.New("author not found")
	ErrAuthorAlreadyExists = errors.New("author already exists")
)
