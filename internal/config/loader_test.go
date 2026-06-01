package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/khanalsaroj/typegenctl/internal/domain"
)

// TestValidate_AppliesTagDefaults guards the historical copy-paste bug where an
// empty backend tag was written to the frontend field, and where defaulting was
// silently dropped because Validate took its argument by value.
func TestValidate_AppliesTagDefaults(t *testing.T) {
	cfg := Default()
	cfg.Services.Frontend.Image.Tag = ""
	cfg.Services.Backend.Image.Tag = ""

	if err := Validate(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Services.Frontend.Image.Tag != "latest" {
		t.Errorf("frontend tag = %q, want latest", cfg.Services.Frontend.Image.Tag)
	}
	if cfg.Services.Backend.Image.Tag != "latest" {
		t.Errorf("backend tag = %q, want latest (regression: default leaked to frontend)", cfg.Services.Backend.Image.Tag)
	}
}

func TestValidate_Errors(t *testing.T) {
	cases := map[string]func(*domain.Config){
		"bad frontend port":      func(c *domain.Config) { c.Services.Frontend.Port.Host = 0 },
		"bad backend port":       func(c *domain.Config) { c.Services.Backend.Port.Host = 99999 },
		"missing backend image":  func(c *domain.Config) { c.Services.Backend.Image.Name = "" },
		"missing frontend image": func(c *domain.Config) { c.Services.Frontend.Image.Name = "" },
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := Default()
			mutate(cfg)
			if err := Validate(cfg); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	if _, err := Load(""); err == nil {
		t.Fatal("expected error for empty path")
	}
}

func TestLoad_NotFound(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.yaml"))
	if !errors.Is(err, domain.ErrConfigNotFound) {
		t.Fatalf("expected ErrConfigNotFound, got %v", err)
	}
}

func TestLoad_MalformedYAML(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bad.yaml")
	// A leading tab is never valid YAML indentation.
	if err := os.WriteFile(p, []byte("\tnot: valid"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(p)
	if !errors.Is(err, domain.InvalidConfig) {
		t.Fatalf("expected InvalidConfig, got %v", err)
	}
}

func TestLoad_FailsValidation(t *testing.T) {
	p := filepath.Join(t.TempDir(), "invalid.yaml")
	body := "" +
		"services:\n" +
		"  frontend:\n" +
		"    image: {name: fe, tag: latest}\n" +
		"    port: {host: 0, container: 80}\n" +
		"    enabled: true\n" +
		"  backend:\n" +
		"    image: {name: be, tag: latest}\n" +
		"    port: {host: 8049, container: 8080}\n" +
		"    enabled: true\n"
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(p)
	if !errors.Is(err, domain.InvalidConfig) {
		t.Fatalf("expected InvalidConfig, got %v", err)
	}
}

func TestLoad_RoundTrip(t *testing.T) {
	p := filepath.Join(t.TempDir(), "typegen.yaml")
	if err := Write(p, Default()); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Services.Backend.Image.Name == "" {
		t.Fatal("expected backend image to be loaded")
	}
}
