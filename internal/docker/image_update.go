package docker

import (
	"bytes"
	"fmt"
	"io"

	dkrImage "github.com/docker/docker/api/types/image"
)

func (c *Client) UpdateImage(imageName string) (bool, error) {

	reader, err := c.Cli.ImagePull(c.ctx, imageName, dkrImage.PullOptions{})
	if err != nil {
		return false, fmt.Errorf("image pull failed: %w", err)
	}
	defer func(reader io.ReadCloser) {
		if reader.Close() != nil {
		}
	}(reader)

	buf, err := io.ReadAll(reader)
	if err != nil {
		return false, fmt.Errorf("failed to read pull output: %w", err)
	}

	if bytes.Contains(buf, []byte("Image is up to date")) {
		return false, nil
	}

	return true, nil
}
