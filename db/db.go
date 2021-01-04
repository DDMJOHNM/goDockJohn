package db

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	Conn *pgx.Conn
}

type User struct {
	Id           int64
	Name         string
	CreatedAt    time.Time
	Password     string
	PasswordHash string
	Token        string
}

func (db *Database) Initialise() (*Database, error) {

	var err error
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	db.Conn = conn
	return db, nil
}

func (db *Database) CreateDb(c echo.Context) error {

	m, err := migrate.New(
		"file://migration",
		os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Migrate unable to connect to database: %v\n", err)
		os.Exit(1)
		return err
	}

	// if err := m.Down(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to migrate down: %v\n", err)
	// 	os.Exit(1)
	// 	return err
	// }

	if err := m.Up(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to migarate up: %v\n", err)
		os.Exit(1)
		return err
	}

	return c.String(http.StatusBadRequest, fmt.Sprintf("Result: %s", "Migrate up success"))

}

func (db *Database) Login(c echo.Context) error {

	u := c.FormValue("username")
	p := c.FormValue("password")

	user := &User{}

	db.Initialise()
	defer db.Conn.Close(context.Background())

	var name string
	var id int64
	var createdAt time.Time
	var passwordhash string

	err := db.Conn.QueryRow(context.Background(), "SELECT users.id, users.name, users.createdat,users.passwordhash FROM users WHERE name=$1", u).Scan(&id, &name, &createdAt, &passwordhash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	user.Id = id
	user.Name = name
	user.Password = ""
	user.PasswordHash = passwordhash
	user.CreatedAt = createdAt

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return err
	}

	user.Token = t

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(p+os.Getenv("PEPPER")))

	switch err {
	case nil:

		db.Initialise()
		defer db.Conn.Close(context.Background())
		_, err = db.Conn.Exec(context.Background(), "UPDATE users set token=$1 WHERE id=($2)", user.Token, &user.Id)

		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})

	case bcrypt.ErrMismatchedHashAndPassword:
		return c.JSON(http.StatusBadRequest, err.Error())
	default:
		return c.JSON(http.StatusOK, user)

	}

}

func (db *Database) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
