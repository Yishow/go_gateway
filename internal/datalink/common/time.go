package common

import (
	"fmt"
	"strings"
	"time"
)

var timeLayouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02 15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05.999999999 -0700 MST",
	"2006-01-02 15:04:05 -0700 MST",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05.999999999",
}

// ParseTimeString parses timestamps stored by databases (SQLite/Postgres).
func ParseTimeString(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("empty time value")
	}
	if plusIndex := strings.Index(value, " m="); plusIndex > 0 {
		value = value[:plusIndex]
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05 -0700 MST", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST m=+0.000000000", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05 -0700 MST m=+0.000000000", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST m=+0.000000000", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05 -0700 MST m=+0.000000000", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05Z07:00", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02T15:04:05.999999999Z07:00", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02T15:04:05Z07:00", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02T15:04:05.999999999", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02T15:04:05", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05.999999999", value); err == nil {
		return parsed, nil
	}
	if parsed, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
		return parsed, nil
	}
	for _, layout := range timeLayouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", value)
}
