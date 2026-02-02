package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
)

func (c *Client) Start(name string) error {

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

func (c *Client) Stop(name string) error {
	state, err := c.ContainerState(name)
	if err != nil {
		return err
	}

	if state != StateRunning {
		return fmt.Errorf("container %q is not running (state: %s)", name, state)
	}

	timeout := 10

	if err := c.Cli.ContainerStop(
		c.ctx,
		name,
		container.StopOptions{Timeout: &timeout},
	); err != nil {
		return fmt.Errorf("failed to stop container %q: %w", name, err)
	}

	return nil
}

func (c *Client) Restart(name string) error {
	state, err := c.ContainerState(name)
	if err != nil {
		return err
	}
	timeout := 10
	switch state {
	case StateRunning, StateExited, StatePaused:
		if err := c.Cli.ContainerRestart(
			c.ctx,
			name,
			container.StopOptions{
				Timeout: &timeout,
			},
		); err != nil {
			return fmt.Errorf("failed to restart container %q: %w", name, err)
		}
		return nil

	default:
		return fmt.Errorf("container %q cannot be restarted from state %q", name, state)
	}
}
