package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/sarojkhanal/typegenctl/internal/domain"
	"gopkg.in/yaml.v3"
)

func Load(path string) (*domain.Config, error) {
	if path == "" {
		return nil, errors.New("config path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrConfigNotFound
		}
		return nil, err
	}

	var cfg domain.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, domain.InvalidConfig
	}

	if err := Validate(cfg); err != nil {
		return nil, domain.InvalidConfig
	}

	return &cfg, nil
}

func Validate(c domain.Config) error {
	if c.Services.Frontend.Port.Host <= 0 || c.Services.Frontend.Port.Host > 65535 {
		return errors.New("frontend port is invalid or missing")
	}

	if c.Services.Backend.Port.Host <= 0 || c.Services.Backend.Port.Host > 65535 {
		return errors.New("backend port is invalid or missing")
	}

	if c.Services.Backend.Image.Name == "" {
		return errors.New("backend image  is required")
	}

	if c.Services.Frontend.Image.Name == "" {
		return errors.New("frontend image is required")
	}

	if c.Services.Frontend.Image.Tag == "" {
		fmt.Println("frontend image tag is missing, using latest")
		c.Services.Frontend.Image.Tag = "latest"
	}

	if c.Services.Backend.Image.Tag == "" {
		fmt.Println("backend image tag is missing, using latest")
		c.Services.Frontend.Image.Tag = "latest"
	}

	return nil
}
