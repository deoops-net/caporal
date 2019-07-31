package main

import (
	"testing"
)

func TestPull(t *testing.T) {
	InitDocker()
	c := Container{
		Repo: "registry.cn-beijing.aliyuncs.com/deoops/dkb-api",
		Tag:  "0.1.13",
	}

	if err := c.Pull(); err != nil {
		t.Fail()
	}
}
