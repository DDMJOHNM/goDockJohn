package handlers

import (
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo"
)

func Upload(db *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {

		return c.String(http.StatusOK, "User successfully created")
	}
}
