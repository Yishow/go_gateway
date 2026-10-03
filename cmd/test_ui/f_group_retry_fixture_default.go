//go:build !f_write_group_fixture

package main

import "go-gateway/internal/datalink/groupdelivery"

const fixtureMaxRetriesEnv = "F_FIXTURE_MAX_RETRIES"

func preflightFixtureRetries() error { return nil }

func groupDeliverySenderConfig() groupdelivery.SenderConfig {
	return groupdelivery.SenderConfig{}
}
