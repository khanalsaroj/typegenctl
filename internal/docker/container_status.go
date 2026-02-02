package docker

import (
	"fmt"
	"time"
)

func (c *Client) StatusService(containerName string) StatusResult {
	res := StatusResult{Image: containerName}

	inspect, err := c.Cli.ContainerInspect(c.ctx, containerName)
	if err != nil {
		res.Exists = false
		res.Running = false
		res.Message = fmt.Sprintf("inspect error: %v", err)
		res.Err = err
		return res
	}

	res.Exists = true
	if inspect.State != nil && inspect.State.Running {
		res.Running = true
		res.Message = "running"
	} else if inspect.State != nil {
		res.Running = false
		res.Message = fmt.Sprintf("stopped (%s)", inspect.State.Status)
	} else {
		res.Running = false
		res.Message = "unknown state"
	}

	if t, err := time.Parse(time.RFC3339Nano, inspect.Created); err == nil {
		res.CreatedAt = t
	}

	for port, bindings := range inspect.NetworkSettings.Ports {
		for _, binding := range bindings {
			res.Ports = append(res.Ports, fmt.Sprintf("%s -> %s", port, binding.HostPort))
		}
	}

	return res
}
