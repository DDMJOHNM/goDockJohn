package main

import (
	"composetest/db"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	db := db.Database{}

	//Handlers
	e.GET("/createdb", db.CreateDb)
	e.GET("/user/:data", db.GetUser)
	e.POST("/login", db.Login)

	r := e.Group("/v1")
	r.Use(middleware.JWT([]byte(os.Getenv("SECRET"))))
	//Register Restricted Routes Here with any handler
	r.GET("/", db.Restricted)
	e.Logger.Fatal(e.Start(":8000"))

}

//todo
//connection pool
//logging
//seed db
//auth restricted routes
//bind struct
