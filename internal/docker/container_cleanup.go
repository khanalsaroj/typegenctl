package docker

import (
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	dkrImage "github.com/docker/docker/api/types/image"
)

func (c *Client) CleanUpImages(imageRef string) (*CleanupResult, error) {
	images, err := c.Cli.ImageList(c.ctx, dkrImage.ListOptions{})
	if err != nil {
		return nil, ErrImageListFailed
	}

	appImages := c.collectMatchingImages(images, imageRef)
	if len(appImages) <= 1 {
		return nil, nil
	}

	sort.Slice(appImages, func(i, j int) bool {
		return appImages[i].CreatedAt > appImages[j].CreatedAt
	})

	res := &CleanupResult{
		KeptImage: appImages[0],
		Removed:   appImages[1:],
	}

	for _, img := range res.Removed {
		ids, err := c.removeContainersForImage(img)
		if err != nil {
			return nil, err
		}
		res.Containers = append(res.Containers, ids...)
	}

	if err := c.removeImages(res.Removed); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) removeContainersForImage(img ImageInfo) ([]string, error) {
	args := filters.NewArgs()
	args.Add("ancestor", img.Reference)

	containers, err := c.Cli.ContainerList(c.ctx, container.ListOptions{
		All:     true,
		Filters: args,
	})
	if err != nil {
		return nil, ErrContainerListFailed
	}

	var removed []string

	for _, cn := range containers {
		if cn.State == "running" {
			continue
		}

		if err := c.Cli.ContainerRemove(c.ctx, cn.ID, container.RemoveOptions{
			Force: true,
		}); err != nil {
			return nil, ErrContainerRemove
		}

		removed = append(removed, cn.ID)
	}

	return removed, nil
}

func (c *Client) collectMatchingImages(images []dkrImage.Summary, imageRef string) []ImageInfo {
	var result []ImageInfo

	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == "<none>:<none>" {
				continue
			}
			if containsAppName(tag, imageRef) {
				result = append(result, ImageInfo{
					ID:        img.ID,
					Reference: tag,
					CreatedAt: img.Created,
				})
			}
		}
	}
	return result
}
func containsAppName(tag, appName string) bool {
	return len(tag) >= len(appName) && tag[:len(appName)] == appName
}
