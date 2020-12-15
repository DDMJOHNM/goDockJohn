package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func main() {

	e := echo.New()

	e.GET("/items/:data", GetItems)

	e.Logger.Fatal(e.Start(":8000"))
}

//http://localhost:8000/items/string?name=john

func GetItems(c echo.Context) error {
	itemName := c.QueryParam("name")
	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("the item name is: %s", itemName))
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Please specify the datatype as String or Json"})
	}
}
