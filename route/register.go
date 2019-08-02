package route

import (
	"net/http"

	"github.com/doops-net/caporal/conf"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func Register(e *echo.Echo) {
	// create container
	// curl -H 'Content-Type:application/json' -d '{"repo": "nginx", "tag": "latest", "name": "mynginx", "opts": {"publish": ["10005:80"]}}' 'localhost:8080/container'
	e.POST("/container", CreateContainer)

	e.DELETE("/container/:name", DeleteContainer)

	e.PUT("/container", UpdateContainer)

	// get container by name
	e.GET("/container/:name", GetContainer)

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
}

func CommonRes(c echo.Context, err *error) {
	if (*err) != nil {
		log.Error((*err).Error())
		_ = c.JSON(http.StatusInternalServerError, &conf.RespMsg{Msg: (*err).Error()})
		return
	}
	_ = c.JSON(http.StatusOK, &conf.RespMsg{Code: conf.S_CONTAINER_SUCCESS})
}
