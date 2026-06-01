package usecase

import (
	"errors"
	"fmt"

	"github.com/khanalsaroj/typegenctl/internal/app"

	"github.com/khanalsaroj/typegenctl/internal/config"
	"github.com/khanalsaroj/typegenctl/internal/docker"
	"github.com/khanalsaroj/typegenctl/internal/result"
)

func Start(opts *app.Options) error {
	if opts.DryRun {
		return startDryRun(opts)
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
	defer func(dkr *docker.Client) {
		err := dkr.Close()
		if err != nil {
		}
	}(dkr)

	selector := app.ServiceSelector{
		Backend:  opts.BackendOnly,
		Frontend: opts.FrontendOnly,
	}

	services := selector.Selected(app.ServicesFromConfig(cfg))

	startService := func(name string) {
		state, err := dkr.ContainerState(name)
		if err != nil {
			r.AddFailure(fmt.Errorf("%s container does not exist", name))
			return
		}
		if state == docker.StateRunning {
			r.AddSuccess(fmt.Sprintf("%s is already running", name))
			return
		}
		if err := dkr.StartContainer(name); err != nil {
			r.AddFailure(fmt.Errorf("failed to start %s: %w", name, err))
			return
		}

		r.AddSuccess(fmt.Sprintf("%s started successfully", name))
	}

	for _, svc := range services {
		startService(svc.ContainerName)
	}

	return app.RenderAndExit(r, opts)
}

func startDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("[dry-run] failed to load config: %w", err)
	}
	selector := app.ServiceSelector{Backend: opts.BackendOnly, Frontend: opts.FrontendOnly}
	services := selector.Selected(app.ServicesFromConfig(cfg))
	for _, svc := range services {
		r.AddSuccess(fmt.Sprintf("[dry-run] would start container %s", svc.ContainerName))
	}
	return app.RenderAndExit(r, opts)
}
