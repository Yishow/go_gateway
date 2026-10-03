//go:build !f_write_group_fixture

package main

import (
	"testing"

	"go-gateway/internal/datalink/groupdelivery"
)

func TestNormalBinaryCapacityEnvironmentIsIgnored(t *testing.T) {
	t.Setenv("F_FIXTURE_GROUP_QUOTA_BYTES", "invalid")
	t.Setenv("F_FIXTURE_GLOBAL_QUOTA_BYTES", "-1")
	if err := preflightFixtureCapacity(); err != nil {
		t.Fatalf("normal capacity preflight returned error: %v", err)
	}
	if got := groupDeliveryQuotaConfig(); got != (groupdelivery.QuotaConfig{GlobalMaxBytes: groupDeliveryQuotaBytes}) {
		t.Fatalf("normal quota config = %#v, want production default", got)
	}
}
