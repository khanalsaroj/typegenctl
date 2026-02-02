package config

import (
	"os"

	"github.com/sarojkhanal/typegenctl/internal/domain"

	"gopkg.in/yaml.v3"
)

func Write(path string, cfg *domain.Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Default() *domain.Config {
	return &domain.Config{
		Services: domain.Services{
			Frontend: domain.Service{
				Image: domain.Image{
					Name: "ghcr.io/khanalsaroj/typegen-ui",
					Tag:  "latest",
				},
				ContainerName: "typegen-frontend",
				Port: domain.Port{
					Host:      7359,
					Container: 80,
				},
				Enabled: true,
			},
			Backend: domain.Service{
				Image: domain.Image{
					Name: "ghcr.io/khanalsaroj/typegen-server",
					Tag:  "latest",
				},
				ContainerName: "typegen-backend",
				Port: domain.Port{
					Host:      8049,
					Container: 8080,
				},
				Enabled: true,
			},
		},
	}
}
