package docker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	dkrImage "github.com/docker/docker/api/types/image"
)

func (c *Client) PullImage(image string) (bool, error) {
	ctx, stop := signal.NotifyContext(c.ctx, os.Interrupt)
	defer stop()

	reader, err := c.Cli.ImagePull(ctx, image, dkrImage.PullOptions{})
	if err != nil {
		return false, err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
		}
	}(reader)
	if err := renderSingleProgressBar(reader); err != nil {
		return false, err
	}
	fmt.Println()
	_, err = io.Copy(io.Discard, reader)
	return true, nil
}

func renderSingleProgressBar(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)

	type layer struct {
		current int64
		total   int64
	}
	layers := make(map[string]*layer)

	for scanner.Scan() {
		var ev PullEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}

		if ev.Error != "" {
			return fmt.Errorf(ev.Error)
		}

		if ev.ID == "" || ev.ProgressDetail.Total == 0 {
			continue
		}

		layers[ev.ID] = &layer{
			current: ev.ProgressDetail.Current,
			total:   ev.ProgressDetail.Total,
		}

		var currentSum, totalSum int64
		for _, l := range layers {
			currentSum += l.current
			totalSum += l.total
		}

		renderBar(currentSum, totalSum)
	}

	return scanner.Err()
}

func renderBar(current, total int64) {
	if total == 0 {
		return
	}

	const width = 40
	percent := float64(current) / float64(total)
	filled := int(percent * width)

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	fmt.Printf("\r[%s] %3.0f%%", bar, percent*100)
}
