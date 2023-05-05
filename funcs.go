package main

import (
	"fmt"
	"html"
	"strings"
	"text/template"
	"time"

	"github.com/SerhiiCho/timeago"
)

var templateFuncs = template.FuncMap{
	"esc":      escapeString,
	"unhtml":   html.UnescapeString,
	"ntos":     shortNumberString,
	"fmtdur":   formatDuration,
	"fmtduryt": formatDurationYt,
	"timeago":  timeAgo,
}

func escapeString(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}

func shortNumberString(n uint64) string {
	if n >= 1000000 {
		if n2 := n / 100000 % 10; n2 != 0 {
			return fmt.Sprint(n/1000000, ",", n2, " млн")
		}
		return fmt.Sprint(n/1000000, " млн")
	}
	if n >= 1000 {
		if n2 := n / 100 % 10; n2 != 0 {
			return fmt.Sprint(n/1000, ",", n2, " тыс")
		}
		return fmt.Sprint(n/1000, " тыс")
	}
	return fmt.Sprint(n)
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h == 0 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}

	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func formatDurationYt(v string) string {
	// Input: PT#H#M#S. Cut first 2 symbols for valid parsing.
	d, err := time.ParseDuration(strings.ToLower(v[2:]))
	if err != nil {
		return ""
	}
	return formatDuration(d)
}

func timeAgo(s string) string {
	t, _ := time.Parse(time.RFC3339, s)
	return timeago.Parse(t.Format("2006-01-02 15:04:05"))
}
