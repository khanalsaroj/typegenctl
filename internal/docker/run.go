package docker

import (
	"fmt"
)

func (c *Client) Run(imageName string, containerName string, containerPort int, hostPort int, backendProxy string, service string) StartResult {

	if err := c.ensureContainer(containerName, imageName, hostPort, containerPort, backendProxy, service); err != nil {
		return StartResult{
			Service: containerName,
			Err:     err,
		}
	}

	state, err := c.ContainerState(containerName)
	if err != nil {
		return StartResult{
			Service: containerName,
			Err:     err,
		}
	}

	switch state {

	case StateRunning:
		return StartResult{
			Service: containerName,
			Message: fmt.Sprintf("%s running on port %d", containerName, hostPort),
		}

	case StateCreated, StateExited, StatePaused, StateDead:
		state, err = c.ContainerState(containerName)
		if err != nil || state != StateRunning {
			return StartResult{
				Service: containerName,
				Err: fmt.Errorf(
					"%s started but is not running (state: %s)",
					containerName, state,
				),
			}
		}

		return StartResult{
			Service: containerName,
			Message: fmt.Sprintf("%s is now running on port %d", containerName, hostPort),
		}

	case StateRestarting:
		return StartResult{
			Service: containerName,
			Message: fmt.Sprintf("%s is restarting", containerName),
		}

	case StateNotFound:
		return StartResult{
			Service: containerName,
			Err:     fmt.Errorf("%s container not found", containerName),
		}

	default:
		return StartResult{
			Service: containerName,
			Err:     fmt.Errorf("%s has unexpected state: %s", containerName, state),
		}
	}
}
