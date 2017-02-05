package web

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"

	storage "store"
)

func ScenariosGet(c echo.Context) error {
	return c.Render(http.StatusOK, "scenario/list.html", nil)
}

func ScenariosGetNew(c echo.Context) error {
	return c.Render(http.StatusOK, "scenario/new.html", nil)
}

func ScenariosGetEdit(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	scenario, found := Storage.GetScenario(id)
	if !found {
		return c.Render(http.StatusNotFound, "", nil)
	}
	return c.Render(http.StatusOK, "scenario/edit.html", scenario)
}

func ScenariosPostEdit(c echo.Context) error {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, found := Storage.GetScenario(id); !found {
		return c.Render(http.StatusNotFound, "NOT FOUND", nil)
	}

	content := c.FormValue("code")
	name := c.FormValue("name")
	active := c.FormValue("active") == "on"

	Storage.AddScenario(storage.Scenario{
		Id:      id,
		Name:    name,
		Active:  active,
		Content: content,
	})
	FlashOk(c, "Сценарий №%d был обновлен", id)

	return c.Redirect(http.StatusSeeOther, "/scenarios")
}

func ScenariosPostNew(c echo.Context) error {

	content := c.FormValue("code")
	name := c.FormValue("name")
	active := c.FormValue("active") == "on"

	if name == "" || content == "" {
		FlashError(c, "Сохранение сценария было отменено: не заполненные данные (содержимое или название)")
		return c.Redirect(http.StatusFound, "/scenarios")
	}

	Storage.AddScenario(storage.Scenario{
		Name:    name,
		Content: content,
		Active:  active,
		History: make([]*storage.ScenarioHistory, 0)},
	)

	FlashOk(c, "Новый сценарий был сохранен")
	return c.Redirect(http.StatusSeeOther, "/scenarios")

}

func ScenarioLogList(c echo.Context) error {

	limitStr := c.QueryParam("limit")
	if limitStr == "" {
		limitStr = "100"
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return c.Render(http.StatusNotFound, "NOT FOUND", nil)
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	scenario, found := Storage.GetScenario(id)
	if !found {
		return c.Render(http.StatusNotFound, "NOT FOUND", nil)
	}
	history := scenario.HistoryList(int(limit))

	return c.Render(http.StatusOK, "scenario/history.html", history)
}

func ScenarioGetLog(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("scenario_id"), 10, 64)
	log_id, _ := strconv.ParseInt(c.Param("log_id"), 10, 64)
	filename := Storage.GetScenarioHistoryLogFile(id, log_id)
	fd, err := os.Open(filename)
	if err != nil {
		return c.String(http.StatusNotFound, "")
	}
	defer fd.Close()
	c.Response().WriteHeader(http.StatusOK)
	io.Copy(c.Response(), fd)
	c.Response().Flush()
	return nil
}
