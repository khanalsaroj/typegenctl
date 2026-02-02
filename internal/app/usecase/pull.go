package usecase

import (
	"errors"
	"fmt"

	"github.com/sarojkhanal/typegenctl/internal/app"

	"github.com/sarojkhanal/typegenctl/internal/config"
	"github.com/sarojkhanal/typegenctl/internal/docker"
	"github.com/sarojkhanal/typegenctl/internal/result"
)

func Pull(opts *app.Options) error {
	if opts.DryRun {
		return pullDryRun(opts)
	}

	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return err
	}
	r := &result.ValidationResult{}

	dkr, err := docker.New()
	if err != nil {
		return errors.New("docker is not running or not reachable")
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

	pullService := func(svc app.Service) {

		imgRes := dkr.CheckImage(svc.Image, svc.ContainerName)

		if imgRes.Status == docker.ImageFound {
			r.AddSuccess(fmt.Sprintf("%s image found locally", svc.Image))
		} else {
			ok, err := dkr.PullImage(svc.Image)
			if err != nil {
				r.AddFailure(fmt.Errorf("failed to pull %s image: %w", svc.Image, err))
				return
			}
			if ok {
				r.AddSuccess(fmt.Sprintf("%s image pulled successfully", svc.Image))
			}
		}
	}

	for _, svc := range services {
		pullService(svc)
	}

	return app.RenderAndExit(r, opts)
}

func pullDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("[dry-run] failed to load config: %w", err)
	}
	selector := app.ServiceSelector{Backend: opts.BackendOnly, Frontend: opts.FrontendOnly}
	services := selector.Selected(app.ServicesFromConfig(cfg))
	for _, svc := range services {
		r.AddSuccess(fmt.Sprintf("[dry-run] would pull %s image (%s)", svc.ContainerName, svc.Image))
	}
	return app.RenderAndExit(r, opts)
}
