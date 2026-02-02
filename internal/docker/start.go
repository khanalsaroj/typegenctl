package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
)

func (c *Client) StartContainer(name string) error {
	state, err := c.ContainerState(name)
	if err != nil {
		return err
	}
	switch state {
	case StateRunning:
		return nil
	case StateCreated, StateExited, StatePaused:
		if err := c.Cli.ContainerStart(c.ctx, name, container.StartOptions{}); err != nil {
			return fmt.Errorf("failed to start container %q: %w", name, err)
		}
		return nil
	default:
		return fmt.Errorf("container %q cannot be started from state %q", name, state)
	}
}
