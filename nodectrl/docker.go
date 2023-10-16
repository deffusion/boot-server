package nodectrl

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"log"
	"strings"
)

type DockerCli struct {
	cli *client.Client
}

func New() (*DockerCli, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.WithMessage(err, "docker-cli.New")
	}
	return &DockerCli{
		cli,
	}, nil
}

func (d *DockerCli) NewContainer(containerConfig *container.Config, hostConfig *container.HostConfig, containerName string) (Container, error) {
	tag := "DockerCli.NewContainer"
	ctx := context.TODO()
	resp, err := d.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		return Container{nil, nil, "", ""}, errors.WithMessage(err, tag)
	}
	if len(resp.Warnings) > 0 {
		warns := strings.Join(resp.Warnings, "\n")
		log.Printf("%s:\n %s", tag, warns)
	}
	return Container{ctx, d, resp.ID, containerName}, nil
}
