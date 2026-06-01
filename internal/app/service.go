package app

import (
	"github.com/khanalsaroj/typegenctl/internal/domain"
)

type ServiceSelector struct {
	Backend  bool
	Frontend bool
}

type Service struct {
	Image         string
	ContainerName string
	ContainerPort int
	HostPort      int
	Enabled       bool
	Service       string
}

func (s ServiceSelector) Selected(services []Service) []Service {
	var selected []Service

	for _, svc := range services {
		if !svc.Enabled {
			continue
		}

		if !s.Backend && !s.Frontend {
			selected = append(selected, svc)
			continue
		}

		if s.Backend && svc.Service == domain.Backend {
			selected = append(selected, svc)
		} else if s.Frontend && svc.Service == domain.Frontend {
			selected = append(selected, svc)
		}

	}

	return selected
}

func ServicesFromConfig(cfg *domain.Config) []Service {
	return []Service{
		{
			Image:         cfg.Services.Frontend.Image.Name + ":" + cfg.Services.Frontend.Image.Tag,
			ContainerPort: cfg.Services.Frontend.Port.Container,
			HostPort:      cfg.Services.Frontend.Port.Host,
			ContainerName: cfg.Services.Frontend.ContainerName,
			Enabled:       cfg.Services.Frontend.Enabled,
			Service:       domain.Frontend,
		},
		{

			Image:         cfg.Services.Backend.Image.Name + ":" + cfg.Services.Backend.Image.Tag,
			ContainerPort: cfg.Services.Backend.Port.Container,
			HostPort:      cfg.Services.Backend.Port.Host,
			ContainerName: cfg.Services.Backend.ContainerName,
			Enabled:       cfg.Services.Backend.Enabled,
			Service:       domain.Backend,
		},
	}
}
