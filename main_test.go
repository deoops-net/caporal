package main

import (
	"fmt"
	"testing"

	"github.com/doops-net/caporal/driver"
)

func TestPull(t *testing.T) {
	driver.InitDocker()
	c := driver.Container{
		Repo: "registry.cn-beijing.aliyuncs.com/deoops/dkb-api",
		Tag:  "0.1.13",
	}

	if err := c.Pull(); err != nil {
		t.Fail()
	}
}

func TestDefer(t *testing.T) {
	funcWithDefer()
}

func echoA(a *int) {
	fmt.Println(*a)
}

func funcWithDefer() (a int) {
	defer echoA(&a)
	a = 1
	return
}
