package docker

import (
	"errors"
	"github.com/khanalsaroj/typegenctl/internal/domain"
	"time"
)

type ContainerState string
type ImageStatus string

var (
	ErrImageListFailed     = errors.New("failed to list docker images")
	ErrContainerListFailed = errors.New("failed to list containers for image")
	ErrContainerRemove     = errors.New("failed to remove container")
	ErrImageRemove         = errors.New("failed to remove image")
)

const (
	ImageFound        ImageStatus = "found"
	ImageNotFound     ImageStatus = "not_found"
	ImageInspectError ImageStatus = "inspect_error"
	ImageTagMismatch  ImageStatus = "tag_mismatch"
)

const (
	StateCreated    ContainerState = "created"
	StateRunning    ContainerState = "running"
	StatePaused     ContainerState = "paused"
	StateRestarting ContainerState = "restarting"
	StateRemoving   ContainerState = "removing"
	StateExited     ContainerState = "exited"
	StateDead       ContainerState = "dead"
	StateNotFound   ContainerState = "not-found"
	StateUnknown    ContainerState = "unknown"
)

type ImageCheckResult struct {
	Exists       bool        `json:"exists"`
	ImageRef     string      `json:"image_ref"`
	Label        string      `json:"label"`
	Status       ImageStatus `json:"status"` // "found", "not_found", "inspect_error", "tag_mismatch"
	RequestedTag string      `json:"requested_tag,omitempty"`
	LocalTags    []string    `json:"local_tags,omitempty"`
	ErrorMessage string      `json:"error_message,omitempty"`
}

type RunService struct {
	Name  string
	RunFn func(*domain.Config) error
	Port  int
}

type ImageInfo struct {
	ID        string
	Reference string
	CreatedAt int64
}

type PullEvent struct {
	ID             string `json:"id"`
	Error          string `json:"error"`
	ProgressDetail struct {
		Current int64 `json:"current"`
		Total   int64 `json:"total"`
	} `json:"progressDetail"`
}

type StartResult struct {
	Service string
	Running bool
	Message string
	Err     error
}

type StatusResult struct {
	Image           string
	Exists          bool
	Running         bool
	Message         string
	Err             error
	CreatedAt       time.Time
	UpdateAvailable bool
	Ports           []string
}
type CleanupResult struct {
	KeptImage  ImageInfo
	Removed    []ImageInfo
	Containers []string
}
