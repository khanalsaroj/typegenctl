package usecase

import (
	"errors"
	"fmt"

	"github.com/sarojkhanal/typegenctl/internal/app"

	"github.com/sarojkhanal/typegenctl/internal/config"
	"github.com/sarojkhanal/typegenctl/internal/docker"
	"github.com/sarojkhanal/typegenctl/internal/result"
)

func Restart(opts *app.Options) error {
	if opts.DryRun {
		return restartDryRun(opts)
	}
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
	defer func(dkr *docker.Client) { _ = dkr.Close() }(dkr)

	selector := app.ServiceSelector{
		Backend:  opts.BackendOnly,
		Frontend: opts.FrontendOnly,
	}

	services := selector.Selected(app.ServicesFromConfig(cfg))

	restartService := func(containerName string) {
		state, err := dkr.ContainerState(containerName)
		if err != nil {
			r.AddFailure(fmt.Errorf("%s container does not exist", containerName))
			return
		}

		if err := dkr.Restart(containerName); err != nil {
			r.AddFailure(fmt.Errorf("failed to restart %s: %w", containerName, err))
			return
		}

		r.AddSuccess(fmt.Sprintf("%s restarted successfully (previous state: %s)", containerName, state))
	}

	for _, svc := range services {

		restartService(svc.ContainerName)
	}

	return app.RenderAndExit(r, opts)
}

func restartDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("[dry-run] failed to load config: %w", err)
	}
	selector := app.ServiceSelector{Backend: opts.BackendOnly, Frontend: opts.FrontendOnly}
	services := selector.Selected(app.ServicesFromConfig(cfg))
	for _, svc := range services {
		r.AddSuccess(fmt.Sprintf("[dry-run] would restart container %s", svc.ContainerName))
	}
	return app.RenderAndExit(r, opts)
}
