package users

import (
	"composetest/db"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (user *User) GetUser(c echo.Context) error {

	// db := &db.Database{}

	// db.Initialise()

	// idd := c.QueryParam("id")

	// var name string
	// var id int64
	// var createdAt time.Time
	// var passwordhash string
	// err := db.Pool.QueryRow(context.Background(), "SELECT users.id, users.name, users.createdat, users.passwordhash FROM users WHERE id=$1", idd).Scan(&id, &name, &createdAt, &passwordhash)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }

	// u := &User{
	// 	Id:           id,
	// 	Name:         name,
	// 	CreatedAt:    createdAt,
	// 	PasswordHash: passwordhash,
	// }

	// defer db.Pool.Close()

	// dataType := c.Param("data")
	// if dataType == "string" {
	// 	return c.String(http.StatusOK, fmt.Sprintf("my user name is : %s my id is: %d  I was created at: %v", name, id, createdAt))
	// } else if dataType == "Json" {
	// 	return c.JSON(http.StatusOK, u)
	// } else {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"error": "Please specify the datatype as String or Json"})
	// }

	return nil

}

func (user *User) CreateUser(c echo.Context) error {

	db := &db.Database{}

	u := c.FormValue("username")
	p := c.FormValue("password")

	// //todo: use bind
	db.Initialise()

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(p+os.Getenv("PEPPER")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Name = u
	user.Password = ""
	user.PasswordHash = string(hashedBytes)
	user.CreatedAt = time.Now()

	// //TODO: check if user exists in DB

	_, err = db.Pool.Exec(context.Background(), "INSERT INTO users(name,createdat,passwordhash) values($1,$2,$3)", &user.Name, &user.CreatedAt, &user.PasswordHash)

	if err != nil {
		return err
	}

	defer db.Pool.Close()

	return c.String(http.StatusOK, "User successfully created")

}
