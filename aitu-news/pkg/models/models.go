package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type Article struct {
	ID              int64
	Title           string
	Text            string
	CategoryId      int64
	PublicationDate time.Time
	ImgPath         string
	FormattedDate   string
}

type User struct {
	ID       int64
	Name     string
	Lastname string
	Email    string
	Password string
	Role     string
}
