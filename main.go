package main

import (
	"composetest/api/users"
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
	//e.Logger.SetLevel(log.DEBUG)
	//e.Logger.Debug("HEllo")
	db := db.Database{}
	users := users.User{}

	//Handlers
	e.GET("/createdb", db.CreateDb)
	e.POST("/createUser", users.CreateUser)
	e.POST("/login", db.Login)

	//Register Restricted Routes Here with any handler
	r := e.Group("/v1")
	r.Use(middleware.JWT([]byte(os.Getenv("SECRET"))))
	r.GET("/", db.Restricted)
	r.GET("/user/:data", users.GetUser)
	e.Logger.Fatal(e.Start(":8000"))

}

//todo
//connection pool
//logging
//seed db
//auth restricted routes
//bind struct
//logout
//aws?? or gcp?? -cloud run??
