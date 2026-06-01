package version

import "testing"

func TestString(t *testing.T) {
	old := Version
	defer func() { Version = old }()

	Version = "v9.9.9"
	if got := String(); got != "v9.9.9" {
		t.Fatalf("String() = %q, want v9.9.9", got)
	}
}
