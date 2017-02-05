package web

import (
	"log"
	"net/http"

	"github.com/labstack/echo"

	auth "auth"
	storage "store"
)

func UsersGet(c echo.Context) error {
	return c.Render(http.StatusOK, "user/list.html", nil)
}

func UsersPost(c echo.Context) error {

	email := c.FormValue("email")
	password := c.FormValue("password")

	if password == "" {
		FlashError(c, "Слишком простой пароль")
		return c.Redirect(http.StatusFound, "/users")
	}

	hash, err := auth.WebHashPassword(password)
	if err != nil {
		FlashError(c, "При создании хэша пароля возникли ошибки: %s", err.Error())
		return c.Redirect(http.StatusFound, "/users")
	}

	log.Printf("[INFO] create new user {email: `%s`, hash: `%s`}", email, hash)
	if err := Storage.AddUser(storage.User{Email: email, Hash: hash}); err != nil {
		FlashError(c, "При добавлении пользователя возникли ошибки: %s", err.Error())
		return c.Redirect(http.StatusNotAcceptable, "/users")
	}

	FlashOk(c, "Пользователь %s был создан", email)
	return c.Redirect(http.StatusSeeOther, "/users")
}

func UserDelete(c echo.Context) error {
	email := c.Param("email")
	log.Printf("[INFO] delete user: {email: `%s`}", email)
	Storage.DelUser(email)
	FlashOk(c, "Пользователь %s был удален.", email)
	return c.Redirect(http.StatusSeeOther, "/users")
}
