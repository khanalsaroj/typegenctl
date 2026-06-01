package usecase

import (
	"errors"
	"fmt"

	"github.com/khanalsaroj/typegenctl/internal/app"

	"github.com/khanalsaroj/typegenctl/internal/config"
	"github.com/khanalsaroj/typegenctl/internal/docker"
	"github.com/khanalsaroj/typegenctl/internal/result"
)

func Update(opts *app.Options) error {
	if opts.DryRun {
		return updateDryRun(opts)
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
		if dkr.Close() != nil {
		}
	}(dkr)

	selector := app.ServiceSelector{
		Backend:  opts.BackendOnly,
		Frontend: opts.FrontendOnly,
	}

	services := selector.Selected(app.ServicesFromConfig(cfg))

	updateService := func(name string, image string) {
		updated, updateErr := dkr.UpdateImage(image)
		if updateErr != nil {
			r.AddFailure(fmt.Errorf("failed to update %s: %w", name, updateErr))
			return
		}
		if !updated {
			r.AddSuccess(fmt.Sprintf("%s image already up to date", name))
			return
		}
		r.AddSuccess(fmt.Sprintf("%s image updated successfully", name))
	}

	for _, svc := range services {
		updateService(svc.ContainerName, svc.Image)
	}

	return app.RenderAndExit(r, opts)
}

func updateDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("[dry-run] failed to load config: %w", err)
	}
	selector := app.ServiceSelector{Backend: opts.BackendOnly, Frontend: opts.FrontendOnly}
	services := selector.Selected(app.ServicesFromConfig(cfg))
	for _, svc := range services {
		r.AddSuccess(fmt.Sprintf("[dry-run] would update image for %s (%s)", svc.ContainerName, svc.Image))
	}
	return app.RenderAndExit(r, opts)
}
