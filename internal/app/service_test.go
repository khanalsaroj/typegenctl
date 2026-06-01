package app

import (
	"testing"

	"github.com/khanalsaroj/typegenctl/internal/domain"
)

func sampleConfig() *domain.Config {
	return &domain.Config{
		Services: domain.Services{
			Frontend: domain.Service{
				Image:         domain.Image{Name: "fe", Tag: "latest"},
				ContainerName: "typegen-frontend",
				Port:          domain.Port{Host: 7359, Container: 80},
				Enabled:       true,
			},
			Backend: domain.Service{
				Image:         domain.Image{Name: "be", Tag: "latest"},
				ContainerName: "typegen-backend",
				Port:          domain.Port{Host: 8049, Container: 8080},
				Enabled:       true,
			},
		},
	}
}

func TestServicesFromConfig(t *testing.T) {
	svcs := ServicesFromConfig(sampleConfig())
	if len(svcs) != 2 {
		t.Fatalf("expected 2 services, got %d", len(svcs))
	}
	if svcs[0].Image != "fe:latest" {
		t.Errorf("frontend image = %q, want fe:latest", svcs[0].Image)
	}
	if svcs[1].Image != "be:latest" {
		t.Errorf("backend image = %q, want be:latest", svcs[1].Image)
	}
	if svcs[1].Service != domain.Backend {
		t.Errorf("second service = %q, want %q", svcs[1].Service, domain.Backend)
	}
}

func TestSelected_NoFlagsSelectsAllEnabled(t *testing.T) {
	got := ServiceSelector{}.Selected(ServicesFromConfig(sampleConfig()))
	if len(got) != 2 {
		t.Fatalf("expected all enabled services, got %d", len(got))
	}
}

func TestSelected_BackendOnly(t *testing.T) {
	got := ServiceSelector{Backend: true}.Selected(ServicesFromConfig(sampleConfig()))
	if len(got) != 1 || got[0].Service != domain.Backend {
		t.Fatalf("expected only backend, got %+v", got)
	}
}

func TestSelected_FrontendOnly(t *testing.T) {
	got := ServiceSelector{Frontend: true}.Selected(ServicesFromConfig(sampleConfig()))
	if len(got) != 1 || got[0].Service != domain.Frontend {
		t.Fatalf("expected only frontend, got %+v", got)
	}
}

func TestSelected_SkipsDisabled(t *testing.T) {
	cfg := sampleConfig()
	cfg.Services.Frontend.Enabled = false
	got := ServiceSelector{}.Selected(ServicesFromConfig(cfg))
	if len(got) != 1 || got[0].Service != domain.Backend {
		t.Fatalf("disabled frontend should be skipped, got %+v", got)
	}
}
