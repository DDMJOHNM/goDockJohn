package users

import (
	"github.com/labstack/echo/v4"
)

func (user *User) CreateUser(c echo.Context) error {

	// db := &db.Database{}

	// u := c.FormValue("username")
	// p := c.FormValue("password")

	// // //todo: use bind
	// db.Initialise()

	// hashedBytes, err := bcrypt.GenerateFromPassword(
	// 	[]byte(p+os.Getenv("PEPPER")), bcrypt.DefaultCost)
	// if err != nil {
	// 	return err
	// }
	// user.Name = u
	// user.Password = ""
	// user.PasswordHash = string(hashedBytes)
	// user.CreatedAt = time.Now()

	// // //TODO: check if user exists in DB

	// _, err = db.Pool.Exec(context.Background(), "INSERT INTO users(name,createdat,passwordhash) values($1,$2,$3)", &user.Name, &user.CreatedAt, &user.PasswordHash)

	// if err != nil {
	// 	return err
	// }

	// defer db.Pool.Close()

	// return c.String(http.StatusOK, "User successfully created")
	return nil
}
