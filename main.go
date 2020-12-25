package main

import (
	"composetest/db"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	db := db.Database{}

	//handlers
	e.GET("/user/:data", db.GetUser)

	e.Logger.Fatal(e.Start(":8000"))

}

//todo - just need to be json responses back end only:
//connection pool
//logging
//golang migrate seed db
//auth
