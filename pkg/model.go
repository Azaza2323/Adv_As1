package models

import (
	"errors"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type News struct {
	ID          int
	Audience    string
	Author      string
	Title       string
	Description string
	Content     string
}
type Depo struct {
	ID    int
	DepId int
	Staff int
}
type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword []byte
	Role           string
}
