package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

var SALT = "caoral-salt"
var AUTH_USER string
var AUTH_PASS string

// TODO
// - [ ] test remote api
// - [ ]

func init() {
	AUTH_PASS = os.Getenv("AUTH_PASS")
	AUTH_USER = os.Getenv("AUTH_USER")

	if len(AUTH_USER) != 0 && len(AUTH_PASS) != 0 {
		token := base64.StdEncoding.EncodeToString(Encrypt([]byte(AUTH_USER+":"+AUTH_PASS), SALT))
		fmt.Println("Authorization turned on, the the token is :", token)
	}

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	InitDocker()
	initServer()
}

type Container struct {
	Repo string           `json:"repo"`
	Tag  string           `json:"tag"`
	Name string           `json:"name"`
	Opts ContainerOptions `json:"opts"`
}

type ContainerOptions struct {
	// Publish equals to -p flag e.g. Publish: {"8080:80", "4431:443"}
	Publish []string `json:"publish"`
}

func (c Container) Start() (err error) {

	dct, err := DockerClient.CreateContainer(docker.CreateContainerOptions{
		Name: c.Name,
		Config: &docker.Config{
			ExposedPorts: c.CreateExposePorts(),
			Image:        c.Repo + ":" + c.Tag,
		},
		HostConfig: &docker.HostConfig{
			PortBindings:    c.CreateBindingPorts(),
			PublishAllPorts: true,
		},
	})
	if err != nil {
		return err
	}

	if err = DockerClient.StartContainer(dct.ID, nil); err != nil {
		return
	}

	return
}

func (c Container) CreateExposePorts() map[docker.Port]struct{} {
	ports := map[docker.Port]struct{}{}

	for _, v := range c.Opts.Publish {
		portMap := strings.Split(v, ":")
		dockerPort := portMap[1]

		ports[docker.Port(dockerPort+"/tcp")] = struct{}{}
		ports[docker.Port(dockerPort+"/udp")] = struct{}{}
	}

	log.Debug("exposed ports:", ports)
	return ports
}

func (c Container) CreateBindingPorts() map[docker.Port][]docker.PortBinding {
	ports := map[docker.Port][]docker.PortBinding{}

	for _, v := range c.Opts.Publish {
		portMap := strings.Split(v, ":")
		dockerPort := portMap[1]
		hostPort := portMap[0]
		ports[docker.Port(dockerPort+"/tcp")] = []docker.PortBinding{docker.PortBinding{HostPort: hostPort, HostIP: "0.0.0.0"}}
		ports[docker.Port(dockerPort+"/udp")] = []docker.PortBinding{docker.PortBinding{HostPort: hostPort, HostIP: "0.0.0.0"}}
	}

	log.Debug("binding ports:", ports)
	return ports
}

func initServer() {
	e := echo.New()
	e.Use(Auth)
	// create container
	// curl -H 'Content-Type:application/json' -d '{"repo": "nginx", "tag": "latest", "name": "mynginx", "opts": {"publish": ["10005:80"]}}' 'localhost:8080/container'
	e.POST("/container", func(c echo.Context) (err error) {
		//createcadvisor()
		ct := Container{}
		if err = c.Bind(&ct); err != nil {
			return err
		}
		log.Debug(ct)
		//var port docker.port = "80"

		if err = ct.Start(); err != nil {
			return
		}

		return
	})

	e.DELETE("/container/:name", func(c echo.Context) (err error) {
		act, err := GetContainerByName(c.Param("name"))

		if err != nil {
			log.Error(err)
			return err
		}

		if reflect.DeepEqual(act, docker.APIContainers{}) {
			log.Debug("no container found")
			return
		}

		if act.State != "running" {
			if err = DockerClient.RemoveContainer(docker.RemoveContainerOptions{
				ID: act.ID,
			}); err != nil {
				log.Error(err)
				return
			}
		}

		if err = DockerClient.StopContainer(act.ID, 0); err != nil {
			log.Error(err)
			return
		}

		return
	})

	e.PUT("/container", func(c echo.Context) (err error) {
		ct := Container{}
		if err = c.Bind(&ct); err != nil {
			return err
		}
		log.Debug(ct)

		act, err := GetContainerByName(ct.Name)
		if err != nil {
			log.Error(err)
			return err
		}
		// no container created before
		// create one
		if reflect.DeepEqual(act, docker.APIContainers{}) {
			log.Debug("no container found")
			if err = ct.Start(); err != nil {
				return
			}
			return
		}

		log.Debug(act)

		if err = DockerClient.StopContainer(act.ID, 0); err != nil {
			return
		}

		if err = DockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID:    act.ID,
			Force: true,
		}); err != nil {
			return
		}

		if err = ct.Start(); err != nil {
			return
		}

		return
	})

	// get container by name
	e.GET("/container/:name", func(c echo.Context) (err error) {
		containerName := c.Param("name")
		container, err := GetContainerByName(containerName)
		if err != nil {
			return err
		}

		log.Debug(container)
		return
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	log.Fatal(e.Start(":8080"))
}

func GetContainerByName(name string) (container docker.APIContainers, err error) {
	cts, err := DockerClient.ListContainers(docker.ListContainersOptions{
		All:     true,
		Size:    false,
		Limit:   0,
		Since:   "",
		Before:  "",
		Filters: nil,
		Context: nil,
	})
	if err != nil {
		return
	}
	//log.Debug(cts)
	for _, v := range cts {
		for _, n := range v.Names {
			if ("/" + name) == n {
				container = v
				log.Debug("got")
				return
			}
		}
	}

	return
}

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		defer func() {
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
		}()

		// not specify auth method skip
		if len(AUTH_PASS) == 0 || len(AUTH_USER) == 0 {
			return nil
		}

		authCode := c.Request().Header.Get("X-AUTH")
		if len(authCode) == 0 {
			return errors.New("no auth field")
		}

		d, _ := base64.StdEncoding.DecodeString(authCode)
		d = Decrypt(d, SALT)
		up := strings.Split(string(d), ":")

		if len(up) != 2 {
			return errors.New("authorize failed")
		}

		if up[0] != AUTH_USER || up[1] != AUTH_PASS {
			return errors.New("authorize failed")
		}

		return nil
	}
}
