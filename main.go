package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/items/:data", GetItems)
	e.Logger.Fatal(e.Start(":8000"))

}

//http://localhost:8000/items/string?name=john

func GetItems(c echo.Context) error {

	connStr := "postgres://postgres:john@localhost/demo_backend?sslmode=verify-full"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stdout, "ERROR &v\n", err)
	}

	var id int64
	var name string
	var createdAt time.Time

	err = db.QueryRow("select id, name, created-at from user where id=$1", 1).Scan(&id, &name, &createdAt)
	if err != nil {
		fmt.Fprintf(os.Stdout, "ERROR &v\n", err)
	}

	fmt.Fprintf(os.Stdout, "ID &v\n", id)
	fmt.Fprintf(os.Stdout, "USER &v\n", name)
	fmt.Fprintf(os.Stdout, "CREATEDAT &v\n", createdAt)

	itemName := c.QueryParam("name")
	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("my user name is: %s", itemName))
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the datatype as String or Json"})
	}
}
