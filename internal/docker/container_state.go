package docker

import "fmt"

func (c *Client) ContainerState(name string) (ContainerState, error) {
	inspect, err := c.Cli.ContainerInspect(c.ctx, name)
	if err != nil {
		return StateNotFound, fmt.Errorf("container not found")
	}

	if inspect.State == nil {
		return StateUnknown, fmt.Errorf("container state unknown")
	}

	return ContainerState(inspect.State.Status), nil
}
