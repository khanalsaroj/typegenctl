package domain

import "errors"

type Config struct {
	Services Services `yaml:"services"`
}

type Services struct {
	Frontend Service `yaml:"frontend"`
	Backend  Service `yaml:"backend"`
}
type Image struct {
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
}

type Service struct {
	Image         Image  `yaml:"image"`
	ContainerName string `yaml:"container_name"`
	Port          Port   `yaml:"port"`
	Enabled       bool   `yaml:"enabled"`
}
type Port struct {
	Host      int `yaml:"host"`
	Container int `yaml:"container"`
}

var ErrConfigNotFound = errors.New("config not found")
var InvalidConfig = errors.New("invalid config")

const (
	Backend  string = "backend"
	Frontend string = "frontend"
)
