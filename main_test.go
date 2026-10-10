package main

import (
	"strings"
	"testing"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		city, wantPath string
	}{
		{"Rome", "https://wttr.in/Rome?"},
		{"New York", "https://wttr.in/New%20York?"},
		{"São Paulo", "https://wttr.in/S%C3%A3o%20Paulo?"},
	}
	for _, tt := range tests {
		got := buildURL(tt.city)
		if !strings.HasPrefix(got, tt.wantPath) {
			t.Errorf("buildURL(%q) = %q, want prefix %q", tt.city, got, tt.wantPath)
		}
		if !strings.Contains(got, "format=") {
			t.Errorf("buildURL(%q) = %q, missing format param", tt.city, got)
		}
	}
}
