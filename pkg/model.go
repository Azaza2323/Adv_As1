package models

import (
	"errors"
)

var ErrNoRecord = errors.New("models: no matching record found")

type News struct {
	ID          int
	Audience    string
	Author      string
	Title       string
	Description string
	Content     string
}
