package web

import (
	"github.com/labstack/echo"

	auth "auth"
)

func AuthSkip(c echo.Context) bool {
	// пропускаем все GET
	if c.Request().Method == "GET" {
		return true
	}
	return false
}

func AuthValidator(username, password string, c echo.Context) bool {
	// если нет пользунов - дефолтный логин/пароль admin/admin
	if len(Storage.Users.List) == 0 && username == "admin" && password == "admin" {
		return true
	}
	if user, found := Storage.GetUser(username); found {
		return auth.WebHashValid(password, user.Hash)
	}
	return false
}
