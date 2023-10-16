package nodectrl

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/pkg/errors"
)

type Container struct {
	ctx  context.Context
	dCli *DockerCli
	id   string
	name string
}

func (c *Container) Start() error {
	return c.dCli.cli.ContainerStart(c.ctx, c.id, types.ContainerStartOptions{})
}

func (c *Container) Stop() error {
	return c.dCli.cli.ContainerStop(c.ctx, c.id, container.StopOptions{})
}

// Wait blocks until the container stops
func (c *Container) Wait() error {
	tag := "Container.Wait:"
	statusCh, errCh := c.dCli.cli.ContainerWait(c.ctx, c.id, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.WithMessage(err, tag)
		}
	case <-statusCh:
	}
	return nil
}
