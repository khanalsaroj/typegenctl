package usecase

import (
	"fmt"
	"strconv"

	"github.com/khanalsaroj/typegenctl/internal/app"
	"github.com/khanalsaroj/typegenctl/internal/browser"
	"github.com/khanalsaroj/typegenctl/internal/config"
	"github.com/khanalsaroj/typegenctl/internal/result"
)

func Dashboard(opts *app.Options) error {
	cfg, err := config.Load(opts.ConfigPath)
	if err != nil {
		return fmt.Errorf(
			"configuration file is missing; initialization is required before this command can be executed: %w",
			err,
		)
	}
	r := &result.ValidationResult{}
	frontendPort := cfg.Services.Frontend.Port.Host
	url := "http://localhost:" + strconv.Itoa(frontendPort)

	fmt.Printf("Dashboard: %s\n", url)
	r.AddSuccess(fmt.Sprintf("Dashboard: %s", url))
	if !browser.CanAutoOpen() {
		return nil
	}

	if err := browser.Open(url); err != nil {
		r.AddFailure(fmt.Errorf("warning: failed to open browser: %w", err))
		r.AddInfo(fmt.Sprintf("Open the URL manually: %s", url))
	}
	return nil
}
