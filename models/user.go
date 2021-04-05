package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type User struct {
	Id           int64
	Name         string
	CreatedAt    time.Time
	Password     string
	PasswordHash string
	Token        string
}

func GetUserByName(db *pgxpool.Pool, username string) (User, error) {

	var user = User{}

	err := db.QueryRow(context.Background(), "SELECT users.id, users.name, users.createdat, users.passwordhash FROM users WHERE name=$1", username).Scan(&user.Id, &user.Name, &user.CreatedAt, &user.PasswordHash)

	if err != nil {
		return user, errors.Wrap(err, "Database scan error")
	}

	defer db.Close()

	return user, nil

}
