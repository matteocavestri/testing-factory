package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: weather <city>")
		os.Exit(1)
	}
	city := strings.Join(os.Args[1:], " ")
	u := "https://wttr.in/" + url.PathEscape(city) + "?format=%l:+%c+%t+(feels+%f),+%C,+humidity+%h,+wind+%w"

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
