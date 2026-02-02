package docker

import (
	"errors"
	dkrImage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/errdefs"
	"strings"
)

func (c *Client) CheckImage(imageRef, label string) *ImageCheckResult {
	result := &ImageCheckResult{
		ImageRef: imageRef,
		Label:    label,
	}

	inspect, err := c.Cli.ImageInspect(c.ctx, imageRef)
	if err != nil {
		var notFound errdefs.ErrNotFound
		if errors.As(err, &notFound) {
			result.Exists = false
			result.Status = ImageNotFound
			return result
		}

		result.Exists = false
		result.Status = ImageInspectError
		result.ErrorMessage = err.Error()
		return result
	}

	result.Exists = true
	result.LocalTags = inspect.RepoTags

	requestedTag := getTag(imageRef)
	result.RequestedTag = requestedTag

	if requestedTag == "" {
		result.Status = ImageFound
		return result
	}

	for _, t := range inspect.RepoTags {
		if getTag(t) == requestedTag {
			result.Status = ImageFound
			return result
		}
	}

	result.Status = ImageTagMismatch
	return result
}

func getTag(image string) string {
	if strings.Contains(image, "@sha256:") {
		return ""
	}

	parts := strings.Split(image, ":")
	if len(parts) < 2 {
		return "latest"
	}

	return parts[len(parts)-1]
}

func (c *Client) removeImages(images []ImageInfo) error {
	for _, img := range images {
		if _, err := c.Cli.ImageRemove(c.ctx, img.ID, dkrImage.RemoveOptions{
			Force:         true,
			PruneChildren: true,
		}); err != nil {
			return ErrImageRemove
		}
	}
	return nil
}
