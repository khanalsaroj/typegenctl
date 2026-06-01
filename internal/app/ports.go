package app

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"syscall"
)

type PortCheckResult struct {
	Name      string
	Port      int
	Available bool
	Err       error
}

func CheckPort(port int, cNName string) PortCheckResult {

	if err := validateContainerName(cNName); err != nil {
		return PortCheckResult{
			Name: cNName,
			Port: port,
			Err:  err,
		}
	}

	if port <= 0 || port > 65535 {
		return PortCheckResult{
			Name: cNName,
			Port: port,
			Err:  fmt.Errorf("invalid %s port: %d (must be 1–65535)", cNName, port),
		}
	}

	available, err := isPortAvailable(port)
	if err != nil {
		return PortCheckResult{
			Name: cNName,
			Port: port,
			Err:  err,
		}
	}

	if !available {
		return PortCheckResult{
			Name:      cNName,
			Port:      port,
			Available: false,
			Err:       fmt.Errorf("%s port %d is already in use", cNName, port),
		}
	}

	return PortCheckResult{
		Name:      cNName,
		Port:      port,
		Available: true,
		Err:       nil,
	}
}

func isPortAvailable(port int) (bool, error) {
	addr := fmt.Sprintf("0.0.0.0:%d", port)

	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		// Address already in use or permission denied
		if errors.Is(err, syscall.EADDRINUSE) {
			return false, nil
		}
		return false, err
	}

	defer func() { _ = ln.Close() }()
	return true, nil
}

var dockerNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]{0,127}$`)

func validateContainerName(cnName string) error {
	if strings.TrimSpace(cnName) == "" {
		return fmt.Errorf("container name cannot be empty")
	}

	if strings.ContainsAny(cnName, " /") {
		return fmt.Errorf("invalid container name '%s': spaces or slashes are not allowed", cnName)
	}

	if !dockerNameRegex.MatchString(cnName) {
		return fmt.Errorf(
			"invalid container name '%s': must match %s (allowed: letters, numbers, ., _, - and must start with alphanumeric)",
			cnName,
			dockerNameRegex.String(),
		)
	}

	return nil
}
