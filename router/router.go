package router

import (
	"fmt"
	"net/http"

	"github.com/AlexKLWS/youtube-audio-stream/consts"
	"github.com/spf13/viper"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Router struct
type Router struct {
	Server    *echo.Echo
	Videos    *echo.Group
	Playlists *echo.Group
}

// New echo router
func New() *Router {
	e := echo.New()

	e.Use(middleware.Recover())

	if viper.GetBool(consts.Debug) {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}))
	}

	e.Static(fmt.Sprintf("/%s", viper.GetString(consts.OutputRoute)), fmt.Sprintf("./%s", viper.GetString(consts.OutputRoute)))
	// Serving the website
	e.Static("/", "client/build")

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "client/build",
		HTML5:  true,
		Browse: false,
	}))

	a := e.Group("/api")

	return &Router{
		Server:    e,
		Videos:    a.Group("/videos"),
		Playlists: a.Group("/playlists"),
	}
}
