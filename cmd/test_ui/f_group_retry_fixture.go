//go:build f_write_group_fixture

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/groupdelivery"
)

const (
	fixtureMaxRetriesEnv = "F_FIXTURE_MAX_RETRIES"
	fixtureMaxRetriesMin = 1
	fixtureMaxRetriesMax = 100
)

func parseFixtureMaxRetries(raw string) (int, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < fixtureMaxRetriesMin || value > fixtureMaxRetriesMax {
		return 0, fmt.Errorf("%s must be an integer from %d through %d", fixtureMaxRetriesEnv, fixtureMaxRetriesMin, fixtureMaxRetriesMax)
	}
	return value, nil
}

func fixtureSenderConfig() (groupdelivery.SenderConfig, error) {
	maxRetries, err := parseFixtureMaxRetries(os.Getenv(fixtureMaxRetriesEnv))
	if err != nil {
		return groupdelivery.SenderConfig{}, err
	}
	return groupdelivery.SenderConfig{MaxRetries: maxRetries}, nil
}

func preflightFixtureRetries() error {
	_, err := fixtureSenderConfig()
	if err != nil {
		return err
	}
	return preflightFixtureWorkerNode()
}

func groupDeliverySenderConfig() groupdelivery.SenderConfig {
	config, err := fixtureSenderConfig()
	if err != nil {
		return groupdelivery.SenderConfig{}
	}
	return config
}
