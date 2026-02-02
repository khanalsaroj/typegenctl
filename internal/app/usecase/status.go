package usecase

import (
	"errors"
	"fmt"
	"github.com/sarojkhanal/typegenctl/internal/app"
	"github.com/sarojkhanal/typegenctl/internal/config"
	"github.com/sarojkhanal/typegenctl/internal/docker"
	"github.com/sarojkhanal/typegenctl/internal/result"
	"strings"
)

func Status(opts *app.Options) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}
	r := &result.ValidationResult{}

	dkr, err := docker.New()
	if err != nil {
		r.AddFailure(errors.New("docker is not running or not reachable"))
		return app.RenderAndExit(r, opts)
	}
	defer func(dkr *docker.Client) {
		if dkr.Close() != nil {
		}
	}(dkr)

	selector := app.ServiceSelector{
		Backend:  opts.BackendOnly,
		Frontend: opts.FrontendOnly,
	}

	services := selector.Selected(app.ServicesFromConfig(cfg))

	statusService := func(containerName string) {
		status := dkr.StatusService(containerName)
		if status.Err != nil {
			r.AddFailure(fmt.Errorf("%s: %v", containerName, status.Err))
		} else {
			msg := fmt.Sprintf("%s: %s", containerName, status.Message)

			if !status.CreatedAt.IsZero() {
				msg += fmt.Sprintf(", Created: %s", status.CreatedAt.Format("2006-01-02 15:04:05"))
			}

			if len(status.Ports) > 0 {
				msg += ", Ports: " + strings.Join(status.Ports, ", ")
			}

			r.AddSuccess(msg)
		}
	}

	for _, svc := range services {
		statusService(svc.ContainerName)
	}
	return app.RenderAndExit(r, opts)
}
