package main

import (
	"testing"

	"github.com/docker/docker/client"
)

func TestPullImage(t *testing.T) {
	executed := false
	err := WithDocker(
		WithImage("alpine:latest"),
		WithClientOpts(client.WithAPIVersionNegotiation()),
		CleanupContainer(),
		CleanupImage(),
		Exec(func(config DockerConfig) (err error) {
			executed = true
			return
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !executed {
		t.Fatal("function should have executed")
	}
}

func TestBuildImage(t *testing.T) {
	executed := false
	err := WithDocker(
		WithDockerfile("./Dockerfile"),
		WithClientOpts(client.WithAPIVersionNegotiation()),
		CleanupContainer(),
		CleanupImage(),
		Exec(func(config DockerConfig) (err error) {
			executed = true
			return
		}),
	)
	if err != nil {
		t.Fatal(err)
	}
	if !executed {
		t.Fatal("function should have executed")
	}
}
