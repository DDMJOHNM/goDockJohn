package users

import (
	"time"
)

type User struct {
	Id           int64
	Name         string
	CreatedAt    time.Time
	Password     string
	PasswordHash string
	Token        string
}
