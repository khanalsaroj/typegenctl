package app

import (
	"net"
	"strings"
	"testing"
)

func TestValidateContainerName(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid simple", "typegen-backend", false},
		{"valid dots underscores", "a.b_c-1", false},
		{"empty", "", true},
		{"whitespace only", "   ", true},
		{"contains space", "bad name", true},
		{"contains slash", "bad/name", true},
		{"leading hyphen", "-bad", true},
		{"too long", strings.Repeat("a", 129), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateContainerName(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("validateContainerName(%q) err = %v, wantErr = %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestCheckPort_InvalidName(t *testing.T) {
	res := CheckPort(8080, "bad name")
	if res.Err == nil {
		t.Fatal("expected error for invalid container name")
	}
	if res.Available {
		t.Fatal("expected Available=false on error")
	}
}

func TestCheckPort_InvalidPort(t *testing.T) {
	for _, p := range []int{0, -1, 70000} {
		if res := CheckPort(p, "typegen-backend"); res.Err == nil {
			t.Fatalf("expected error for port %d", p)
		}
	}
}

func TestCheckPort_InUse(t *testing.T) {
	ln, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		t.Skipf("cannot bind: %v", err)
	}
	defer func() { _ = ln.Close() }()
	port := ln.Addr().(*net.TCPAddr).Port

	res := CheckPort(port, "typegen-backend")
	if res.Available {
		t.Fatalf("expected port %d to be reported in use", port)
	}
	if res.Err == nil {
		t.Fatal("expected error when port in use")
	}
}

func TestCheckPort_Available(t *testing.T) {
	ln, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		t.Skipf("cannot bind: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close() // release it

	res := CheckPort(port, "typegen-backend")
	if res.Err != nil {
		t.Fatalf("unexpected error: %v", res.Err)
	}
	if !res.Available {
		t.Fatalf("expected port %d to be available", port)
	}
}
