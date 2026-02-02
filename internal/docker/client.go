package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/client"
)

type Client struct {
	Cli *client.Client
	ctx context.Context
}

func New() (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	if _, err := cli.Ping(ctx); err != nil {
		_ = cli.Close()
		return nil, err
	}

	return &Client{
		Cli: cli,
		ctx: context.Background(),
	}, nil
}

func (c *Client) Close() error {
	if c.Cli == nil {
		return nil
	}
	return c.Cli.Close()
}
