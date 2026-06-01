package usecase

import (
	"errors"
	"fmt"
	"github.com/khanalsaroj/typegenctl/internal/app"
	"github.com/khanalsaroj/typegenctl/internal/domain"

	"github.com/khanalsaroj/typegenctl/internal/config"
	"github.com/khanalsaroj/typegenctl/internal/result"
)

func Initialization(opts *app.Options) error {
	if opts.DryRun {
		return initializationDryRun(opts)
	}

	r := &result.ValidationResult{}
	_, err := config.Load(opts.ConfigPath)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConfigNotFound):
			cfg := buildConfig(opts)
			if err := config.Write(opts.ConfigPath, cfg); err != nil {
				return fmt.Errorf("failed to write default config: %w", err)
			}
			r.AddSuccess(fmt.Sprintf("typegen.yaml created at %s", opts.ConfigPath))

		case errors.Is(err, domain.InvalidConfig):
			r.AddFailure(fmt.Errorf(
				"invalid config file %s\nfix the file or delete it and rerun init\noriginal error: %w",
				opts.ConfigPath,
				err,
			))
		default:
			return fmt.Errorf("failed to load config: %w", err)
		}
	} else {
		if !opts.Force {
			r.AddSuccess("typegen.yaml already exists (use --force to overwrite port for backend and frontend services)")
		} else {
			cfg := buildConfig(opts)
			if err := config.Write(opts.ConfigPath, cfg); err != nil {
				return fmt.Errorf("failed to overwrite config: %w", err)
			}
			r.AddSuccess("typegen.yaml overwritten using")
		}
	}

	return app.RenderAndExit(r, opts)
}

func buildConfig(opts *app.Options) *domain.Config {
	cfg := config.Default()

	if opts.FrontendPort > 0 {
		cfg.Services.Frontend.Port.Host = opts.FrontendPort
	}
	if opts.BackendPort > 0 {
		cfg.Services.Backend.Port.Host = opts.BackendPort
	}

	return cfg
}

func initializationDryRun(opts *app.Options) error {
	r := &result.ValidationResult{}
	_, err := config.Load(opts.ConfigPath)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConfigNotFound):
			_ = buildConfig(opts)
			r.AddSuccess("[dry-run] default config.yaml would be created")
		case errors.Is(err, domain.InvalidConfig):
			r.AddFailure(fmt.Errorf(
				"[dry-run] invalid config file %s\\nfix the file or delete it and rerun init\\noriginal error: %w",
				opts.ConfigPath,
				err,
			))
		default:
			return fmt.Errorf("[dry-run] failed to load config: %w", err)
		}
	} else {
		if !opts.Force {
			r.AddSuccess("[dry-run] typegen.yaml already exists (use --force to overwrite port for backend and frontend services)")
		} else {
			_ = buildConfig(opts)
			r.AddSuccess("[dry-run] typegen.yaml would be overwritten using provided options")
		}
	}
	return app.RenderAndExit(r, opts)
}
