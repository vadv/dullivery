package web

import (
	"net/http"

	"github.com/labstack/echo"
)

func ConfigGet(c echo.Context) error {
	return c.Render(http.StatusOK, "config/show.html", nil)
}

func ConfigPost(c echo.Context) error {
	newConfig := c.FormValue("config")
	if Storage.Config != newConfig {
		Storage.SaveConfig(newConfig)
		FlashOk(c, "Конфигурация была сохранена")
		return c.Render(http.StatusOK, "config/show.html", nil)
	} else {
		FlashError(c, "Сохранение конфигурации было отменено: нечего обновлять")
		return c.Render(http.StatusOK, "config/show.html", nil)
	}
}
