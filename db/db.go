package db

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

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

//http://localhost:8000/user/json?id=1

func (db *Database) GetUser(c echo.Context) error {

	db.Initialise()

	idd := c.QueryParam("id")

	var name string
	var id int64
	var createdAt time.Time

	err := db.Conn.QueryRow(context.Background(), "SELECT id, name, createdat FROM users WHERE id=$1", idd).Scan(&id, &name, &createdAt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	u := &User{
		Id:        id,
		Name:      name,
		CreatedAt: createdAt,
	}

	defer db.Conn.Close(context.Background())

	// // fmt.Print(u)
	dataType := c.Param("data")
	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("my user name is : %s my id is: %d  I was created at: %v", name, id, createdAt))
	} else if dataType == "Json" {
		return c.JSON(http.StatusOK, u)
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the datatype as String or Json"})
	}
	//return c.String(http.StatusBadRequest, fmt.Sprintf("User: %s", "hello"))
}
