package main

import (
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
	flag.Parse()

	if *units != "metric" && *units != "imperial" {
		fmt.Fprintf(os.Stderr, "error: invalid units %q (must be 'metric' or 'imperial')\n", *units)
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: weather <city> [--units metric|imperial]")
		os.Exit(1)
	}
	city := strings.Join(flag.Args(), " ")
	u := buildURL(city, *units)

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
	fmt.Println(strings.TrimSpace(string(body)))
}

func buildURL(city, units string) string {
	var u string
	switch units {
	case "imperial":
		u = "?u"
	case "metric":
		u = "?m"
	}
	return "https://wttr.in/" + url.PathEscape(city) + "?format=" + u + "%l:+%c+%t+(feels+%f),+%C,+humidity+%h,+wind+%w"
}
