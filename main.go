package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	units := flag.String("units", "metric", "units: metric or imperial")
	flag.StringVar(units, "u", "metric", "units: metric or imperial")
	forecast := flag.Bool("f", false, "show a 3-day forecast")
	flag.BoolVar(forecast, "forecast", false, "show a 3-day forecast")
	flag.Parse()

	if *units != "metric" && *units != "imperial" {
		fmt.Fprintf(os.Stderr, "error: invalid units %q (must be 'metric' or 'imperial')\n", *units)
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: weather <city> [--units metric|imperial] [--forecast]")
		os.Exit(1)
	}
	city := strings.Join(flag.Args(), " ")
	u := buildURL(city, *units, *forecast)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: city %q not found (%s)\n", city, resp.Status)
		os.Exit(1)
	}

	var data WeatherJSON
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot parse weather data: %v\n", err)
		os.Exit(1)
	}
	if !data.HasWeather() {
		fmt.Fprintf(os.Stderr, "error: city %q not found\n", city)
		os.Exit(1)
	}

	if *forecast {
		fmt.Print(FormatForecast(data.CurrentDesc(), data.Days()))
	} else {
		fmt.Println(strings.TrimSpace(string(body)))
	}
}

// WeatherJSON mirrors the wttr.in ?format=j1 response.
type WeatherJSON struct {
	Current  []Condition `json:"current_condition"`
	Weather  json.RawMessage `json:"weather"`
}

// Condition is the current-conditions entry from wttr.in. wttr.in returns the
// value either as a bare string or as an object with a "value" field, so this
// custom UnmarshalJSON accepts both shapes.
type Condition struct {
	v string
}

// UnmarshalJSON accepts a bare string or a {"value": "..."} object.
func (c *Condition) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '{' {
		var obj struct {
			Value *string `json:"value"`
		}
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		if obj.Value != nil {
			c.v = *obj.Value
		}
		return nil
	}
	c.v = string(data)
	return nil
}

// String returns the condition value without any JSON quoting.
func (c Condition) String() string {
	return strings.Trim(c.v, "\"")
}

// HasWeather reports whether the wttr.in response contained a weather field.
func (w WeatherJSON) HasWeather() bool {
	return len(w.Weather) > 0
}

// Days parses the weather field into a slice of days.
func (w WeatherJSON) Days() []Day {
	var days []Day
	if err := json.Unmarshal(w.Weather, &days); err != nil {
		return nil
	}
	return days
}

// CurrentDesc is the description of the current conditions (first element of
// the current_condition array in the wttr.in response).
func (w WeatherJSON) CurrentDesc() string {
	if len(w.Current) > 0 {
		return w.Current[0].String()
	}
	return ""
}

// Day is a single day of forecast from wttr.in.
type Day struct {
	Date       string `json:"date"`
	MaxTempC   string `json:"maxtempC"`
	MinTempC   string `json:"mintempC"`
	DescMaxDay string `json:"descMaxDay"`
}

// FormatForecast prints the current conditions line followed by one line for
// each of the next three days. It performs no HTTP request.
func FormatForecast(current string, days []Day) string {
	var b strings.Builder
	if current != "" {
		b.WriteString(current)
		b.WriteString("\n")
	}
	for i := 0; i < 3 && i < len(days); i++ {
		d := days[i]
		line := fmt.Sprintf("%s  %s  +%s°C / +%s°C\n", d.Date, d.DescMaxDay, d.MaxTempC, d.MinTempC)
		b.WriteString(line)
	}
	return b.String()
}

func buildURL(city, units string, forecast bool) string {
	var u string
	switch units {
	case "imperial":
		u = "?u"
	case "metric":
		u = "?m"
	}
	if forecast {
		u += "&format=j1"
	}
	return "https://wttr.in/" + url.PathEscape(city) + "?format=" + u
}
