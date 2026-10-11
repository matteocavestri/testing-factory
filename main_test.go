package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadFixture(t *testing.T, name string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("reading testdata/%s: %v", name, err)
	}
	return data
}

func TestFormatForecast(t *testing.T) {
	got := FormatForecast("Rome Value", []Day{
		{"Sat 12 Oct", "22", "14", "Partly cloudy"},
		{"Sun 13 Oct", "25", "13", "Sunny"},
		{"Mon 14 Oct", "26", "15", "Sunny"},
		{"Tue 15 Oct", "24", "12", "Light rain"},
	})
	want := `Rome Value
Sat 12 Oct  Partly cloudy  +22°C / +14°C
Sun 13 Oct  Sunny  +25°C / +13°C
Mon 14 Oct  Sunny  +26°C / +15°C
`
	if got != want {
		t.Errorf("FormatForecast() =\n%q\nwant\n%q", got, want)
	}
}

func TestFormatForecastTruncatesToThreeDays(t *testing.T) {
	got := FormatForecast("", []Day{
		{"Mon", "1", "0", "a"},
		{"Tue", "2", "1", "b"},
		{"Wed", "3", "2", "c"},
		{"Thu", "4", "3", "d"},
	})
	if strings.Count(got, "\n") != 3 {
		t.Errorf("expected exactly 3 day lines, got:\n%q", got)
	}
	if !strings.Contains(got, "Wed") {
		t.Errorf("expected oldest kept day, got:\n%q", got)
	}
	if strings.Contains(got, "Thu") {
		t.Errorf("expected newest day truncated, got:\n%q", got)
	}
}

func TestParseForecastJSON(t *testing.T) {
	data := loadFixture(t, "forecast.json")
	var w WeatherJSON
	if err := json.Unmarshal(data, &w); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !w.HasWeather() {
		t.Error("HasWeather = false, want true")
	}
	if desc := w.CurrentDesc(); desc != "Rome Value" {
		t.Errorf("CurrentDesc() = %q, want %q", desc, "Rome Value")
	}
	days := w.Days()
	if len(days) != 4 {
		t.Fatalf("Days length = %d, want 4", len(days))
	}
	if days[0].Date != "Sat 12 Oct" {
		t.Errorf("Days[0].Date = %q, want %q", days[0].Date, "Sat 12 Oct")
	}
}

func TestParseForecastMissingWeather(t *testing.T) {
	data := []byte(`{"current_condition": ["x"]}`)
	var w WeatherJSON
	if err := json.Unmarshal(data, &w); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.HasWeather() {
		t.Error("HasWeather = true, want false")
	}
}

func TestParseForecastMalformed(t *testing.T) {
	data := []byte(`{ this is not valid json`)
	var w WeatherJSON
	if err := json.Unmarshal(data, &w); err == nil {
		t.Error("expected error on malformed JSON")
	}
}

func TestBuildURLForecast(t *testing.T) {
	cases := []struct {
		city string
		want string
	}{
		{"Rome", "https://wttr.in/Rome?format=?m&format=j1"},
		{"New York", "https://wttr.in/New%20York?format=?m&format=j1"},
	}
	for _, c := range cases {
		if got := buildURL(c.city, "metric", true); !strings.HasPrefix(got, c.want) {
			t.Errorf("buildURL(%q) forecast = %q, want prefix %q", c.city, got, c.want)
		}
	}
}

func TestBuildURLNoForecast(t *testing.T) {
	if got := buildURL("Rome", "metric", false); got != "https://wttr.in/Rome?format=?m" {
		t.Errorf("buildURL metric no forecast = %q", got)
	}
	if got := buildURL("Rome", "imperial", false); got != "https://wttr.in/Rome?format=?u" {
		t.Errorf("buildURL imperial no forecast = %q", got)
	}
}
