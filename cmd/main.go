package main

import (
	"composetest/bindings"
	handlers "composetest/handlers"
	"composetest/models"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

func main() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)

	e.Validator = new(bindings.Validator)

	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	e.Logger.Debug(dbpool)

	e.GET("/health-check", handlers.HealthCheck)
	e.POST("/login", handlers.Login(dbpool))
	e.POST("/createUser", handlers.CreateUser(dbpool))
	e.GET("/createdb", handlers.CreateDb(dbpool))

	var signingKey = []byte(os.Getenv("SECRET"))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(models.SigningContextKey, signingKey)
			return next(c)
		}
	})

	v1 := e.Group("/v1")

	users := v1.Group("/user", middleware.JWT(signingKey))
	users.Use(middleware.JWT(signingKey))
	users.GET("/:id", handlers.GetUserByID(dbpool))

	//uploads := v1.Group("/upload", middleware.JWT(signingKey))
	//uploads.Use(middleware.JWT(signingKey))
	e.POST("upload", handlers.Upload(dbpool))

	e.Logger.Fatal(e.Start(":8000"))

}

//cloud run
