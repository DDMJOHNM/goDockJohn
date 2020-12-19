package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/jackc/pgx/v4"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/items/:data", GetItems)
	e.Logger.Fatal(e.Start(":8000"))

}

//http://localhost:8000/items/string?name=john

func GetItems(c echo.Context) error {

	var err error
	conn, err := pgx.Connect(context.Background(), "postgresql://postgres:john@composetest_database_1/demo_backend?sslmode=disable")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	var name string
	var id int64
	var createdAt time.Time

	err = conn.QueryRow(context.Background(), "select id, name, createdat from users where id=$1", 1).Scan(&id, &name, &createdAt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	//itemName := c.QueryParam("name")
	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("my user is : %s my id is: %d  I was created at: %v", name, id, createdAt))
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the datatype as String or Json"})
	}
}

//todo:
//connection pool
//env var
//hot reload
//logging
//golang migrate seed
