package usecase

import (
	"errors"
	"fmt"

	"github.com/khanalsaroj/typegenctl/internal/app"
	"github.com/khanalsaroj/typegenctl/internal/config"
	"github.com/khanalsaroj/typegenctl/internal/docker"
	"github.com/khanalsaroj/typegenctl/internal/result"
)

func Check(opts *app.Options) error {

	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf(
			"configuration file is missing; initialization is required before this command can be executed: %w",
			err,
		)
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

	r.AddSuccess("docker daemon is running")
	selector := app.ServiceSelector{
		Backend:  opts.BackendOnly,
		Frontend: opts.FrontendOnly,
	}

	services := selector.Selected(app.ServicesFromConfig(cfg))

	checkService := func(imageName string, hostPort int, containerName string) {
		res := app.CheckPort(hostPort, containerName)
		if res.Err != nil {
			r.AddFailure(res.Err)
			return
		}
		if res.Available {
			r.AddSuccess(fmt.Sprintf("%s port %d is available", res.Name, res.Port))
		} else {
			r.AddFailure(fmt.Errorf("%s port %d is already in use", res.Name, res.Port))
		}

		resDkr := dkr.CheckImage(imageName, containerName)
		switch resDkr.Status {
		case docker.ImageFound:
			r.AddInfo(fmt.Sprintf("%s exists locally: %s", resDkr.Label, resDkr.ImageRef))
		case docker.ImageNotFound:
			r.AddFailure(fmt.Errorf("%s not found locally: %s", resDkr.Label, resDkr.ImageRef))
		case docker.ImageInspectError:
			r.AddInfo(fmt.Sprintf("%s inspect error: %s - %s", resDkr.Label, resDkr.ImageRef, resDkr.ErrorMessage))
		case docker.ImageTagMismatch:
			r.AddInfo(fmt.Sprintf("%s tag mismatch: requested=%s, local=%v", resDkr.Label, resDkr.RequestedTag, resDkr.LocalTags))
		}
	}

	for _, svc := range services {
		checkService(svc.Image, svc.HostPort, svc.ContainerName)
	}

	return app.RenderAndExit(r, opts)
}
