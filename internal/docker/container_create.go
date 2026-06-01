package docker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/khanalsaroj/typegenctl/internal/domain"
)

const (
	// defaultNetworkName is the bridge network all TypeGen containers join so
	// the frontend can reach the backend by container name.
	defaultNetworkName = "bridge-net"
	// backendDataTarget is the in-container mount point for the backend volume.
	backendDataTarget = "/app/data"
	// envAppEnv is injected into the backend container.
	envAppEnv = "APP_ENV=production"
	// envAPIUpstreamKey tells the frontend which container to proxy /api to.
	envAPIUpstreamKey = "API_UPSTREAM"
)

func (c *Client) ensureContainer(
	name string,
	image string,
	hostPort int,
	containerPort int,
	backendProxy string,
	serviceType string,
) error {
	networkName := defaultNetworkName

	if err := ensureNetwork(c.Cli, c.ctx, networkName); err != nil {
		return err
	}

	inspect, err := c.Cli.ContainerInspect(c.ctx, name)
	if err == nil && inspect.State != nil && inspect.State.Running {
		return nil
	}

	portKey := nat.Port(fmt.Sprintf("%d/tcp", containerPort))

	var env []string
	switch serviceType {
	case domain.Backend:
		env = append(env, envAppEnv)
	case domain.Frontend:
		if backendProxy != "" {
			env = append(env, envAPIUpstreamKey+"="+backendProxy)
		}
	default:
		return fmt.Errorf("unknown serviceType: %s", serviceType)
	}

	var mounts []mount.Mount
	if serviceType == domain.Backend {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: name,
			Target: backendDataTarget,
		})
	}

	resp, err := c.Cli.ContainerCreate(
		c.ctx,
		&container.Config{
			Image:        image,
			ExposedPorts: nat.PortSet{portKey: struct{}{}},
			Env:          env,
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				portKey: {{
					HostIP:   "0.0.0.0",
					HostPort: strconv.Itoa(hostPort),
				}},
			},
			Mounts: mounts,
		},
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				networkName: {},
			},
		},
		nil,
		name,
	)
	if err != nil {
		return err
	}

	return c.Cli.ContainerStart(c.ctx, resp.ID, container.StartOptions{})
}

func ensureNetwork(cli *client.Client, ctx context.Context, networkName string) error {
	_, err := cli.NetworkCreate(ctx, networkName, network.CreateOptions{
		Driver: "bridge",
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	return nil
}
