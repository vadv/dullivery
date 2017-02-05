package web

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
)

func writeCookie(c echo.Context, name, value string) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = url.QueryEscape(value)
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func writeFlashCookie(c echo.Context, typ, format string, a ...interface{}) {
	writeCookie(c, "flash.text", fmt.Sprintf(format, a...))
	writeCookie(c, "flash.type", typ)
}

func FlashError(c echo.Context, format string, a ...interface{}) {
	writeFlashCookie(c, "danger", format, a...)
}

func FlashOk(c echo.Context, format string, a ...interface{}) {
	writeFlashCookie(c, "success", format, a...)
}
