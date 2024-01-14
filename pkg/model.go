package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type News struct {
	ID          int       `json:"id"`
	Audience    string    `json:"audience"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Created     time.Time `json:"created"`
}
