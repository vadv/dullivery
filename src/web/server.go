package web

import (
	"os"
	"path/filepath"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	storage "store"
)

var Storage *storage.Storage

func Serve(address, staticDir string, fd *os.File) (err error) {

	e := echo.New()

	// load box
	Storage = storage.Box

	// load templates
	e.Renderer = LoadTemplates(staticDir)

	//e.Use(middleware.Recover())
	if fd != nil {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Output: fd}))
	} else {
		e.Use(middleware.Logger())
	}
	e.Use(middleware.BasicAuthWithConfig(
		middleware.BasicAuthConfig{
			Skipper:   AuthSkip,
			Validator: AuthValidator,
		}))

	if fd != nil {
		e.Logger.SetOutput(fd)
	}

	e.Static("/static", filepath.Join(staticDir, "static"))
	e.Static("/fonts", filepath.Join(staticDir, "fonts"))

	e.GET("/", RootGet)

	e.GET("/config", ConfigGet)
	e.POST("/config", ConfigPost)

	e.GET("/users", UsersGet)
	e.POST("/users/new", UsersPost)
	e.POST("/users/delete/:email", UserDelete)

	e.GET("/scenarios", ScenariosGet)
	e.GET("/scenario/new", ScenariosGetNew)
	e.POST("/scenarios/new", ScenariosPostNew)
	e.GET("/scenario/edit/:id", ScenariosGetEdit)
	e.POST("/scenario/edit/:id", ScenariosPostEdit)
	e.GET("/scenario/history/:id", ScenarioLogList)
	e.GET("/scenario/log_file/:scenario_id/:log_id", ScenarioGetLog)

	e.GET("/tasks", TasksGet)
	e.GET("/task/new", TaskNewGet)
	e.POST("/task/new", TaskNewPost)
	e.GET("/task/edit/:id", TaskEditGet)
	e.POST("/task/edit/:id", TaskEditPost)
	e.GET("/task_log/:id", TaskLogGet)

	return e.Start(address)

}

func RootGet(c echo.Context) error {
	return c.Redirect(302, "/scenarios")
}
