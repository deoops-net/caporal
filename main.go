package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/doops-net/caporal/route"

	"github.com/doops-net/caporal/conf"

	"github.com/doops-net/caporal/driver"

	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

var SALT = "caoral-salt"

// TODO
// - [ ] test remote api
// - [ ]

func init() {
	conf.AUTH_PASS = os.Getenv("AUTH_PASSWORD")
	conf.AUTH_USER = os.Getenv("AUTH_USER")
	if len(conf.AUTH_USER) != 0 && len(conf.AUTH_PASS) != 0 {
		token := base64.StdEncoding.EncodeToString(Encrypt([]byte(conf.AUTH_USER+":"+conf.AUTH_PASS), SALT))
		fmt.Println("Authorization turned on, the the token is :", token)
	}

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	driver.InitDocker()
	startServer()
}

func startServer() {
	e := echo.New()
	e.Use(Auth)
	route.Register(e)

	log.Fatal(e.Start(":8080"))
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
			// cipher method may panic
			if panicErr := recover(); panicErr != nil {
				c.String(http.StatusUnauthorized, "authorize failed, wrong crypto method")
				return
			}

			if err != nil {
				c.String(http.StatusUnauthorized, err.Error())
				return
			}

			if err = next(c); err != nil {
				c.Error(err)
			}

			log.Debug("auth success")
		}()

		// not specify auth method skip
		if len(conf.AUTH_PASS) == 0 || len(conf.AUTH_USER) == 0 {
			return nil
		}

		authCode := c.Request().Header.Get("X-AUTH")
		if len(authCode) == 0 {
			return errors.New("no auth field")
		}

		d, _ := base64.StdEncoding.DecodeString(authCode)
		d = Decrypt(d, SALT)
		log.Debug(string(d))
		up := strings.Split(string(d), ":")

		if len(up) != 2 {
			return errors.New("authorize failed")
		}

		if up[0] != conf.AUTH_USER || up[1] != conf.AUTH_PASS {
			return errors.New("authorize failed")
		}

		return nil
	}
}
