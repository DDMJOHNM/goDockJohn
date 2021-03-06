package handlers

import (
	"composetest/bindings"
	"composetest/models"
	"composetest/old_code/api/users"
	"composetest/renderings"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *pgxpool.Pool) echo.HandlerFunc {

	return func(c echo.Context) error {

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

		user, err := models.GetUserByName(db, lr.Username)
		if err != nil {
			resp.Success = false
			resp.Message = "Username or Password incorrect"
			return c.JSON(http.StatusUnauthorized, resp)
		}

		//defer db.Close()

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = user.Name
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		t, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			return err
		}

		resp.Token = t

		//c.Logger().Debug(user)
		//c.Logger().Debug(resp)

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.PasswordHash),
			[]byte(lr.Password+os.Getenv("PEPPER")))

		switch err {
		case nil:

			//defer db.Close()
			_, err = db.Exec(context.Background(), "UPDATE users set token=$1 WHERE id=($2)", resp.Token, &user.Id)

			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, map[string]string{
				"token": t,
			})

		case bcrypt.ErrMismatchedHashAndPassword:
			return c.JSON(http.StatusBadRequest, err.Error())
		default:
			return c.JSON(http.StatusOK, user)

		}
	}

}

func CreateUser(db *pgxpool.Pool) echo.HandlerFunc {

	return func(c echo.Context) error {

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

		_, err = db.Exec(context.Background(), "INSERT INTO users(name,createdat,passwordhash) values($1,$2,$3)", user.Name, &user.CreatedAt, &user.PasswordHash)

		if err != nil {
			return err
		}

		//defer db.Close()

		return c.String(http.StatusOK, "User successfully created")

	}
}

func GetUserByID(db *pgxpool.Pool) echo.HandlerFunc {

	return func(c echo.Context) error {
		user := new(users.User)

		idd, err := strconv.ParseInt(c.FormValue("id"), 10, 64)
		if err != nil {
			return err
		}

		err = db.QueryRow(context.Background(), "SELECT id, name, createdAt FROM public.users WHERE id=$1", idd).Scan(&user.Id, &user.Name, &user.CreatedAt)
		if err != nil {
			return err
		}

		c.Logger().Debug(user)

		return c.JSON(http.StatusOK, user)
	}
}

func CreateDb(db *pgxpool.Pool) echo.HandlerFunc {

	//todo add db seed for users here

	return func(c echo.Context) error {

		m, err := migrate.New(
			"file://migration",
			os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Migrate unable to connect to database: %v\n", err)
			os.Exit(1)
			return err
		}

		if err := m.Down(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to migrate down: %v\n", err)
			os.Exit(1)
			return err
		}

		if err := m.Up(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to migarate up: %v\n", err)
			os.Exit(1)
			return err
		}

		return c.String(http.StatusBadRequest, fmt.Sprintf("Result: %s", "Migrate up success"))
	}

}

func HealthCheck(c echo.Context) error {
	resp := renderings.HealthCheckResponse{
		Message: "Everything is good!",
	}
	return c.JSON(http.StatusOK, resp)
}
