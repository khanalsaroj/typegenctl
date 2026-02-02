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
)

func (c *Client) ensureContainer(
	name string,
	image string,
	hostPort int,
	containerPort int,
	backendProxy string,
	serviceType string,
) error {
	networkName := "bridge-net"

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
	case "backend":
		env = append(env, "APP_ENV=production")
	case "frontend":
		if backendProxy != "" {
			env = append(env, "API_UPSTREAM="+backendProxy)
		}
	default:
		return fmt.Errorf("unknown serviceType: %s", serviceType)
	}

	var mounts []mount.Mount
	if serviceType == "backend" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: name,
			Target: "/app/data",
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
