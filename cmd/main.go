package main

import (
	"composetest/bindings"
	"composetest/handlers"
	"composetest/models"
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

// type jwtCustomClaims struct {
// 	Name  string `json:"name"`
// 	Admin bool   `json:"admin"`
// 	jwt.StandardClaims
// }

// func handler(c echo.Context) error {
// 	return c.String(http.StatusOK, "Hello World")
// }

func main() {

	e := echo.New()

	e.GET("/health-check", handlers.HealthCheck)
	//g := e.Group("/v1")
	e.POST("/login", handlers.Login)
	e.POST("/createUser", handlers.CreateUser)

	e.Logger.Fatal(e.Start(":8000"))
	e.Validator = new(bindings.Validator)

	dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	//defer dbpool.Close()

	// db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	e.Logger.Infof("%v", "Unable to connect to database:", err)
	// }

	e.Logger.Infof("%v", dbpool)

	// err = db.Ping()
	// if err != nil {
	// 	e.Logger.Infof("%v", "Unable to connect to database:", err)
	// }

	// dbpool, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	// 	os.Exit(1)
	// }

	//defer dbpool.Close()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(models.DBContextKey, dbpool)
			fmt.Fprintf(os.Stdout, "Unable to connect to database: %v\n", dbpool)
			return next(c)
		}
	})

	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		c.Set(models.DBContextKey, dbpool)
	// 		return next(c)
	// 	}
	// })

	//g.POST("/logout", handlers.logout)

	// e := echo.New()
	// e.Use(middleware.Logger())
	// //e.Logger.SetLevel(log.DEBUG)
	// //e.Logger.Debug("HEllo")
	// db := db.Database{}
	// db.Initialise()

	// //users := users.User{}

	// //Handlers
	// //e.GET("/createdb", db.CreateDb)
	// //e.POST("/createUser", users.CreateUser)
	// //e.POST("/login", db.Login)

	// //Register Restricted Routes Here with any handler
	// r := e.Group("/v1")
	// r.Use(middleware.JWT([]byte(os.Getenv("SECRET"))))
	// r.GET("/", db.Restricted)
	// r.GET("/user/:data", db.GetUser)
	// e.Logger.Fatal(e.Start(":8000"))

}

//todo
//connection pool
//lift out db init call to main.go
//logging
//seed db
//auth restricted routes
//bind struct
//logout
//aws?? or gcp?? -cloud run??
