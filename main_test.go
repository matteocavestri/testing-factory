package main

import (
	"strings"
	"testing"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		city, units, wantPath string
	}{
		{"Rome", "metric", "https://wttr.in/Rome?format=?m"},
		{"New York", "metric", "https://wttr.in/New%20York?format=?m"},
		{"São Paulo", "metric", "https://wttr.in/S%C3%A3o%20Paulo?format=?m"},
		{"Rome", "imperial", "https://wttr.in/Rome?format=?u"},
		{"New York", "imperial", "https://wttr.in/New%20York?format=?u"},
		{"São Paulo", "imperial", "https://wttr.in/S%C3%A3o%20Paulo?format=?u"},
	}
	for _, tt := range tests {
		got := buildURL(tt.city, tt.units)
		if !strings.HasPrefix(got, tt.wantPath) {
			t.Errorf("buildURL(%q) = %q, want prefix %q", tt.city, got, tt.wantPath)
		}
		if !strings.Contains(got, "format=") {
			t.Errorf("buildURL(%q) = %q, missing format param", tt.city, got)
		}
	}
}
