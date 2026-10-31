package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

func TestProductionWordFormatRestartAndSQLRollback(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			f := newWordProductionFixture(t, protocol)
			rule, p, want := f.createRule(t, protocol, 0, schema.DataTypeUint64, "DCBA")
			require.NoError(t, f.services.runtime.Start(t.Context()))
			before := f.services.scheduler.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, before, 1)
			require.Equal(t, want, before[0].Value)
			// A real SQL persistence fault must not leave new decoding intent installed.
			_, err := f.db.Exec(`CREATE TRIGGER word_fail_rule_point BEFORE UPDATE ON points BEGIN SELECT RAISE(ABORT,'word-test-write-fault'); END`)
			require.NoError(t, err)
			format := "ABCD"
			_, _, err = f.services.sourceRule.UpdateWithRuntimeReconcile(t.Context(), rule.ID, sourcerule.UpdateRuleRequest{DataFormat: &format})
			require.Error(t, err)
			_, err = f.db.Exec(`DROP TRIGGER word_fail_rule_point`)
			require.NoError(t, err)
			stored, err := f.services.point.GetByID(t.Context(), p.ID)
			require.NoError(t, err)
			require.Equal(t, p.RevisionID, stored.RevisionID)
			require.Equal(t, p.FormatMode, stored.FormatMode)
			storedRule, err := f.services.sourceRule.GetByID(t.Context(), rule.ID)
			require.NoError(t, err)
			require.Equal(t, rule.DataFormat, storedRule.DataFormat)
			// Restart constructs new services/scheduler from the exact persisted state.
			require.NoError(t, f.services.runtime.Stop(t.Context()))
			_, err = f.db.Exec(`UPDATE points SET last_read_at=NULL,last_value=NULL,last_error='' WHERE id=?`, p.ID)
			require.NoError(t, err)
			manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
			t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
			restarted, err := wireGatewayServices(f.db, manager)
			require.NoError(t, err)
			t.Cleanup(func() { require.NoError(t, restarted.runtime.Stop(context.Background())) })
			restartedPoint, e := restarted.point.GetByID(t.Context(), p.ID)
			require.NoError(t, e)
			require.Equal(t, p.RevisionID, restartedPoint.RevisionID, "no configuration edit occurred during restart")
			require.NoError(t, restarted.runtime.Start(t.Context()))
			after := restarted.scheduler.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, after, 1)
			require.Equal(t, want, after[0].Value)
			require.Equal(t, before[0].ConfigFingerprint, after[0].ConfigFingerprint)
			require.NotEqual(t, before[0].InstallationScope, after[0].InstallationScope)
			require.False(t, restarted.scheduler.AcceptsInstallationSample(p.ID, before[0].InstallationLease, before[0].InstallationScope))
			raw, err := json.Marshal(want)
			require.NoError(t, err)
			require.Eventually(t, func() bool {
				got, e := restarted.point.GetByID(t.Context(), p.ID)
				return e == nil && got.LastValue != nil && *got.LastValue == string(raw)
			}, 5*time.Second, 20*time.Millisecond)
		})
	}
}

func TestProductionWordFormatZeroAndShortResponse(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			f := newWordProductionFixture(t, protocol)
			_, p, _ := f.createRule(t, protocol, 0, schema.DataTypeUint32, "ABCD")
			f.plc.set(0, []uint16{0, 0})
			require.NoError(t, f.services.runtime.Start(t.Context()))
			good := f.services.scheduler.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, good, 1)
			require.Equal(t, uint32(0), good[0].Value)
			require.Equal(t, schema.QualityGood, good[0].Quality)
			f.plc.mu.Lock()
			f.plc.bad = true
			f.plc.mu.Unlock()
			bad := f.services.scheduler.PollNowContext(t.Context(), []string{p.ID})
			require.Len(t, bad, 1)
			require.NotEqual(t, schema.QualityGood, bad[0].Quality)
			require.Nil(t, bad[0].Value)
			require.NotEmpty(t, bad[0].Error)
		})
	}
}
