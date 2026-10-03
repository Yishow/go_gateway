//go:build !f_write_group_fixture

package main

import "go-gateway/internal/datalink/groupdelivery"

func preflightFixtureCapacity() error { return nil }

func groupDeliveryQuotaConfig() groupdelivery.QuotaConfig {
	return groupdelivery.QuotaConfig{GlobalMaxBytes: groupDeliveryQuotaBytes}
}
