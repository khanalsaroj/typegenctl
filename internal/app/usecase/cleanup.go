package usecase

import (
	"errors"
	"fmt"

	"github.com/sarojkhanal/typegenctl/internal/app"

	"github.com/sarojkhanal/typegenctl/internal/config"
	"github.com/sarojkhanal/typegenctl/internal/docker"
	"github.com/sarojkhanal/typegenctl/internal/result"
)

func CleanUp(opts *app.Options) error {
	if opts.DryRun {
		return cleanupDryRun(opts)
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

	cleanupService := func(svc app.Service) {
		res, err := dkr.CleanUpImages(svc.Image)
		if err != nil {
			r.AddFailure(fmt.Errorf("cleanup failed for %s: %w", svc.ContainerName, err))
			return
		}
		if res == nil {
			return
		}
		r.AddSuccess(fmt.Sprintf("kept latest %s image: %s", svc.ContainerName, res.KeptImage.Reference))

		for _, img := range res.Removed {
			r.AddSuccess(fmt.Sprintf("removed old %s image: %s", svc.ContainerName, img.Reference))
		}
		for _, id := range res.Containers {
			r.AddSuccess(fmt.Sprintf("removed container %s", id[:12]))
		}
	}

	for _, svc := range services {
		cleanupService(svc)
	}

	return app.RenderAndExit(r, opts)
}

func cleanupDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("[dry-run] failed to load config: %w", err)
	}
	selector := app.ServiceSelector{Backend: opts.BackendOnly, Frontend: opts.FrontendOnly}
	services := selector.Selected(app.ServicesFromConfig(cfg))
	for _, svc := range services {
		r.AddSuccess(fmt.Sprintf("[dry-run] would clean up Docker artifacts for %s (image %s): keep latest image, remove older images and stopped containers", svc.ContainerName, svc.Image))
	}
	r.AddInfo("[dry-run] Running containers would not be stopped or removed")
	return app.RenderAndExit(r, opts)
}
