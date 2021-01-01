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
)

type Database struct {
	Conn *pgx.Conn
}

type User struct {
	Id        int64
	Name      string
	CreatedAt time.Time
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

	if err := m.Down(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to migrate down: %v\n", err)
		os.Exit(1)
		return err
	}

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

	//todo: query db here login check and save current session in user table
	//create user
	//use basic auth to send toekn with request see restricted group.

	if u != "john" || p != "secret" {
		return echo.ErrUnauthorized
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "john"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return err
	}

	//on successful login response is returned.

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func (db *Database) Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func (db *Database) GetUser(c echo.Context) error {

	db.Initialise()

	idd := c.QueryParam("id")

	var name string
	var id int64
	var createdAt time.Time

	err := db.Conn.QueryRow(context.Background(), "SELECT users.id, users.name, users.createdat FROM users WHERE id=$1", idd).Scan(&id, &name, &createdAt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	//TODO:Bind User Struct
	u := &User{
		Id:        id,
		Name:      name,
		CreatedAt: createdAt,
	}

	defer db.Conn.Close(context.Background())

	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("my user name is : %s my id is: %d  I was created at: %v", name, id, createdAt))
	} else if dataType == "Json" {
		return c.JSON(http.StatusOK, u)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the datatype as String or Json"})
	}

}
