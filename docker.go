package main

import docker "github.com/fsouza/go-dockerclient"

var DockerClient *docker.Client

func InitDocker() {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	DockerClient = client
}
