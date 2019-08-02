package route

import (
	"reflect"

	"github.com/doops-net/caporal/driver"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func CreateContainer(c echo.Context) (err error) {
	defer CommonRes(c, &err)
	//createcadvisor()
	ct := driver.Container{}
	if err = c.Bind(&ct); err != nil {
		return
	}

	if err = ct.Start(); err != nil {
		return
	}

	return
}

func DeleteContainer(c echo.Context) (err error) {
	defer CommonRes(c, &err)

	act, err := driver.GetContainerByName(c.Param("name"))
	if err != nil {
		return err
	}

	if reflect.DeepEqual(act, docker.APIContainers{}) {
		log.Debug("no container found")
		return
	}

	if act.State != "running" {
		if err = driver.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID: act.ID,
		}); err != nil {
			return
		}
	}

	if err = driver.DockerClient.StopContainer(act.ID, 0); err != nil {
		return
	}

	return
}
func UpdateContainer(c echo.Context) (err error) {
	log.Debug("alksd")
	defer CommonRes(c, &err)
	ct := driver.Container{}
	if err = c.Bind(&ct); err != nil {
		return err
	}

	if err = ct.Pull(); err != nil {
		return
	}

	// TODO implement it for ct
	act, err := driver.GetContainerByName(ct.Name)
	if err != nil {
		return err
	}
	// no container found created before
	// create one
	if reflect.DeepEqual(act, docker.APIContainers{}) {
		log.Debug("no container found")
		if err = ct.Start(); err != nil {
			return
		}
		return
	}

	if err = driver.DockerClient.StopContainer(act.ID, 0); err != nil {
		if _, ok := err.(*docker.ContainerNotRunning); ok {
			log.Debug("container is stopped")
		} else {
			return
		}
	}

	if err = driver.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:    act.ID,
		Force: true,
	}); err != nil {
		return
	}

	if err = ct.Start(); err != nil {
		return
	}

	return
}
func GetContainer(c echo.Context) (err error) {
	defer CommonRes(c, &err)
	containerName := c.Param("name")
	container, err := driver.GetContainerByName(containerName)
	if err != nil {
		return err
	}

	log.Debug(container)
	return
}
