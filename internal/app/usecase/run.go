package usecase

import (
	"errors"
	"fmt"

	"github.com/sarojkhanal/typegenctl/internal/app"
	"github.com/sarojkhanal/typegenctl/internal/config"
	"github.com/sarojkhanal/typegenctl/internal/docker"
	"github.com/sarojkhanal/typegenctl/internal/result"
)

func Run(opts *app.Options) error {
	if opts.DryRun {
		return runDryRun(opts)
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
		_ = dkr.Close()
	}(dkr)

	selector := app.ServiceSelector{
		Backend:  opts.BackendOnly,
		Frontend: opts.FrontendOnly,
	}

	backendProxy := cfg.Services.Backend.ContainerName

	services := selector.Selected(app.ServicesFromConfig(cfg))
	results := make([]docker.StartResult, 0, len(services))
	allRunning := true
	runService := func(imageName string, containerName string, containerPort int, hostPort int, service string) {

		runRes := dkr.Run(imageName, containerName, containerPort, hostPort, backendProxy, service)

		results = append(results, runRes)

		if runRes.Err != nil {
			allRunning = false
			r.AddFailure(runRes.Err)
		}
		r.AddSuccess(runRes.Message)
	}

	for _, svc := range services {
		runService(svc.Image, svc.ContainerName, svc.ContainerPort, svc.HostPort, svc.Service)
	}

	if allRunning {
		fmt.Println()
		fmt.Println("You can access the frontend at:")
		fmt.Printf("http://localhost:%d\n\n", cfg.Services.Frontend.Port.Host)
	}

	return app.RenderAndExit(r, opts)
}

func runDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf("[dry-run] failed to load config: %w", err)
	}
	selector := app.ServiceSelector{Backend: opts.BackendOnly, Frontend: opts.FrontendOnly}
	services := selector.Selected(app.ServicesFromConfig(cfg))
	for _, svc := range services {
		r.AddSuccess(fmt.Sprintf("[dry-run] would run container %s from image %s (container port %d -> host %d)", svc.ContainerName, svc.Image, svc.ContainerPort, svc.HostPort))
	}
	if (!opts.BackendOnly) && (opts.FrontendOnly || (!opts.BackendOnly && !opts.FrontendOnly)) {
		r.AddInfo(fmt.Sprintf("[dry-run] You would access the frontend at: http://localhost:%d", cfg.Services.Frontend.Port.Host))
	}
	return app.RenderAndExit(r, opts)
}
