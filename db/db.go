package db

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx"
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

func Initialise() (Database, error) {

	db := Database{}
	var err error

	pgxConfig, err := pgx.ParseConnectionString(os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse connection string: %v\n", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(pgxConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	db.Conn = conn

	defer db.Conn.Close()

	return db, nil
}

//http://localhost:8000/item/json?id=1

func (db *Database) GetItems(c echo.Context) error {

	idd := c.QueryParam("id")

	var name string
	var id int64
	var createdAt time.Time

	err := db.Conn.QueryRow("select id, name, createdat from users where id=$1", idd).Scan(&id, &name, &createdAt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	u := &User{
		Id:        id,
		Name:      name,
		CreatedAt: createdAt,
	}

	defer db.Conn.Close()

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("my user name is : %s my id is: %d  I was created at: %v", name, id, createdAt))
	} else if dataType == "Json" {
		return c.String(http.StatusOK, fmt.Sprintf("User: %v", u))
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the datatype as String or Json"})
	}
}
