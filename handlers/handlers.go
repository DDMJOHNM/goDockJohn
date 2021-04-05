package handlers

import (
	"composetest/bindings"
	"composetest/models"
	"composetest/renderings"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
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

	db := c.Get(models.DBContextKey).(*pgxpool.Pool)

	user, err := models.GetUserByName(db, lr.Username)
	if err != nil {
		resp.Success = false
		resp.Message = "Username or Password incorrect"
		return c.JSON(http.StatusUnauthorized, resp)
	}

	return c.JSON(http.StatusOK, user)

}

func HealthCheck(c echo.Context) error {
	resp := renderings.HealthCheckResponse{
		Message: "Everything is good!",
	}
	return c.JSON(http.StatusOK, resp)
}
