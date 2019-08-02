package driver

import (
	"fmt"
	"os"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/labstack/gommon/log"
)

var DockerClient *docker.Client

// InitDocker initial a client
func InitDocker() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	DockerClient = client
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
	// Network equals to --network flag
	Network string       `json:"network"`
	Mount   []HostVolume `json:"mount"`
}

type HostVolume struct {
	Bind string `json:"bind"`
	Type string `json:"type"`
}

func (c Container) Pull() (err error) {

	if err = DockerClient.PullImage(docker.PullImageOptions{
		Repository:    c.Repo,
		Tag:           c.Tag,
		OutputStream:  nil,
		RawJSONStream: false,
	}, GenRegistryAuth()); err != nil {
		log.Error(err)
		return
	}

	return
}

func (c Container) Start() (err error) {

	fmt.Println("aisjdoajsod")
	fmt.Println(c.Opts.Network)
	log.Debug(c.Opts.Network)
	dct, err := DockerClient.CreateContainer(docker.CreateContainerOptions{
		Name: c.Name,
		Config: &docker.Config{
			ExposedPorts: c.CreateExposePorts(),
			Image:        c.Repo + ":" + c.Tag,
		},
		HostConfig: &docker.HostConfig{
			PortBindings:    c.CreateBindingPorts(),
			PublishAllPorts: true,
			NetworkMode:     c.Opts.Network,
			//Binds:           c.Opts.Binds,
			Mounts: c.CreateHostMounts(),
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

func (c Container) CreateHostMounts() []docker.HostMount {
	hostMounts := []docker.HostMount{}
	if c.Opts.Mount != nil {
		for _, b := range c.Opts.Mount {
			data := strings.Split(b.Bind, ":")
			if len(data) != 2 {
				continue
			}
			hm := docker.HostMount{
				Target:   data[1],
				Source:   data[0],
				ReadOnly: false,
				Type:     b.Type,
			}
			hostMounts = append(hostMounts, hm)
		}
	}

	return hostMounts
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

func GenRegistryAuth() docker.AuthConfiguration {
	auth := docker.AuthConfiguration{}
	if nonPrivate := os.Getenv("NOT_PRIVATE"); nonPrivate == "true" {
		return auth
	}

	// try reg  auth env  first
	reg_user := os.Getenv("REG_USER")
	reg_passwd := os.Getenv("REG_PASSWORD")
	if len(reg_passwd) != 0 && len(reg_user) != 0 {
		auth.Username = reg_user
		auth.Password = reg_passwd
		return auth
	}
	auth_user := os.Getenv("AUTH_USER")
	auth_passwd := os.Getenv("AUTH_PASSWORD")
	if len(auth_user) != 0 && len(auth_passwd) != 0 {
		auth.Username = auth_user
		auth.Password = auth_passwd
		return auth
	}

	// try user autn env

	return auth
}
