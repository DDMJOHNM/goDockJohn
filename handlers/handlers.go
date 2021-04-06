package handlers

import (
	"composetest/bindings"
	"composetest/models"
	"composetest/old_code/api/users"
	"composetest/renderings"
	"context"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func Login(c echo.Context) error {

	resp := renderings.LoginResponse{}
	lr := new(bindings.LoginRequest)
	if err := c.Bind(lr); err != nil {
		resp.Success = false
		resp.Message = "unable to send bin request for login"
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := lr.Validate(); err != nil {
		resp.Success = false
		resp.Message = err.Error()
		return c.JSON(http.StatusBadRequest, resp)
	}

	//db := c.Get(models.DBContextKey).(*pgxpool.Pool)

	// user, err := models.GetUserByName(db, lr.Username)
	// if err != nil {
	// 	resp.Success = false
	// 	resp.Message = "Username or Password incorrect"
	// 	return c.JSON(http.StatusUnauthorized, resp)
	// }

	return c.JSON(http.StatusOK, c.Get(models.DBContextKey))

}

func CreateUser(c echo.Context) error {

	user := new(users.User)

	resp := renderings.LoginResponse{}
	lr := new(bindings.LoginRequest)
	if err := c.Bind(lr); err != nil {
		resp.Success = false
		resp.Message = "unable to send bin request for login"
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := lr.Validate(); err != nil {
		resp.Success = false
		resp.Message = err.Error()
		return c.JSON(http.StatusBadRequest, resp)
	}

	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(lr.Password+os.Getenv("PEPPER")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Name = lr.Username
	user.Password = ""
	user.PasswordHash = string(hashedBytes)
	user.CreatedAt = time.Now()

	// // //TODO: check if user exists in DB
	db := c.Get(models.DBContextKey).(*pgxpool.Pool)

	_, err = db.Exec(context.Background(), "INSERT INTO users(name,createdat,passwordhash) values($1,$2,$3)", user.Name, &user.CreatedAt, &user.PasswordHash)

	if err != nil {
		return err
	}

	defer db.Close()

	return c.String(http.StatusOK, "User successfully created")
	//return nil
}

func HealthCheck(c echo.Context) error {
	resp := renderings.HealthCheckResponse{
		Message: "Everything is good!",
	}
	return c.JSON(http.StatusOK, resp)
}
