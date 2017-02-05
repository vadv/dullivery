package web

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"

	storage "store"
)

func TasksGet(c echo.Context) error {
	return c.Render(http.StatusOK, "task/list.html", nil)
}

func TaskNewGet(c echo.Context) error {
	return c.Render(http.StatusOK, "task/new.html", nil)
}

func TaskEditGet(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return c.Render(http.StatusOK, "task/edit.html", Storage.GetTask(id))
}

func TaskNewPost(c echo.Context) error {

	content := c.FormValue("content")
	date := c.FormValue("date")

	if content == "" || date == "" {
		FlashError(c, "При создании задания не были указаны: содержимое или время")
		return c.Redirect(http.StatusFound, "/tasks")
	}

	t, err := time.Parse("02/01/2006 15:04", date)
	if err != nil {
		FlashError(c, "При создании задания возникли ошибки: %s", err.Error())
		return c.Redirect(http.StatusFound, "/tasks")
	}

	Storage.AddTask(&storage.Task{
		Content: content,
		StartAt: storage.UnixTime(t.Unix()),
	})
	FlashOk(c, "Задание создано")
	return c.Redirect(http.StatusOK, "/tasks")
}

func TaskLogGet(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	filename := Storage.GetTaskLog(id)
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

func TaskEditPost(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	task := Storage.GetTask(id)
	if task == nil {
		return c.Render(http.StatusNotAcceptable, "", nil)
	}
	content := c.FormValue("content")
	date := c.FormValue("date")
	if date != "" {
		t, err := time.Parse("02/01/2006 15:04", date)
		if err != nil {
			FlashError(c, "При сохранении задания %d возникли ошибки: %s", id, err.Error())
			return c.Redirect(http.StatusFound, "/tasks")
		}
		task.StartAt = storage.UnixTime(t.Unix())
	}
	task.Content = content
	Storage.AddTask(task)
	FlashOk(c, "Задание №%d было обновлено", id)
	return c.Redirect(http.StatusSeeOther, "/tasks")
}
