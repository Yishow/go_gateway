package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/recordingplan"
	"go-gateway/internal/datalink/workspace"
)

func TestSQLServerSuccessorLiveCatalogDamageRefusesActivation(t *testing.T) {
	sqlServerSuccessorFixtures(t) // Missing/unsafe real fixtures must mark the entrypoint NOT_RUN.
	for _, mode := range []workspace.WriteGroupStorageMode{"plain", "native_monthly"} {
		for _, fault := range []string{"owner", "receipt_owner", "disabled_index"} {
			t.Run(string(mode)+"/"+fault, func(t *testing.T) {
				r := newSQLServerSuccessorRig(t, mode, false, true)
				f := r.activationAPI()
				review := f.postActivation(t, f.activationRequest(t))
				saved, err := r.f.groups.Get(t.Context(), r.prepared.Group.ID)
				require.NoError(t, err)
				table, owner, ok := workspace.LongStorageIdentity(saved.Group)
				require.True(t, ok)
				if fault == "receipt_owner" {
					// Never fault a shared receipt table inherited from earlier runs.
					var tokenID string
					require.NoError(t, r.f.db.QueryRow(`SELECT token FROM managed_schema_preview_tokens WHERE plan_id=?`, r.prepared.Group.ID).Scan(&tokenID))
					token, err := recordingplan.NewSQLRepository(r.f.db).GetPreviewToken(t.Context(), tokenID)
					require.NoError(t, err)
					created := false
					for _, candidate := range token.Tables {
						created = created || candidate.Name == dbtarget.EffectReceiptTable && candidate.Action == recordingplan.SchemaTableActionCreate
					}
					require.True(t, created, "receipt fault needs a fresh test-created receipt table; preexisting shared storage must stay untouched")
					table, owner = dbtarget.EffectReceiptTable, "_gw_effect_receipts_v1"
				}
				qualified := sqlServerQualified(r.targetNamespace, table)
				if fault == "disabled_index" {
					var index string
					require.NoError(t, r.targetDB.QueryRowContext(t.Context(), `SELECT name FROM sys.indexes WHERE object_id=OBJECT_ID(@p1) AND is_primary_key=0 AND type>0`, r.targetNamespace+"."+table).Scan(&index))
					require.Regexp(t, `^idx_gw_long_[0-9a-f]{16}_src_time$`, index)
					t.Cleanup(func() {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						_, err := r.targetDB.ExecContext(ctx, "ALTER INDEX ["+index+"] ON "+qualified+" REBUILD")
						require.NoError(t, err)
					})
					_, err = r.targetDB.ExecContext(t.Context(), "ALTER INDEX ["+index+"] ON "+qualified+" DISABLE")
					require.NoError(t, err)
				} else {
					var constraint, definition string
					require.NoError(t, r.targetDB.QueryRowContext(t.Context(), `SELECT d.name,d.definition FROM sys.default_constraints d JOIN sys.columns c ON c.default_object_id=d.object_id WHERE d.parent_object_id=OBJECT_ID(@p1) AND c.name=@p2`, r.targetNamespace+"."+table, owner).Scan(&constraint, &definition))
					require.Regexp(t, `^[A-Za-z0-9_]+$`, constraint)
					require.Contains(t, definition, "gateway")
					t.Cleanup(func() {
						ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
						defer cancel()
						_, err := r.targetDB.ExecContext(ctx, "ALTER TABLE "+qualified+" DROP CONSTRAINT ["+constraint+"]; ALTER TABLE "+qualified+" ADD CONSTRAINT ["+constraint+"] DEFAULT "+definition+" FOR ["+owner+"]")
						require.NoError(t, err)
					})
					_, err = r.targetDB.ExecContext(t.Context(), "ALTER TABLE "+qualified+" DROP CONSTRAINT ["+constraint+"]; ALTER TABLE "+qualified+" ADD CONSTRAINT ["+constraint+"] DEFAULT N'foreign' FOR ["+owner+"]")
					require.NoError(t, err)
				}
				before, external := sqlServerSuccessorLocalState(t, r.f.db), r.targetState(t)
				ready := f.postActivation(t, map[string]any{"action": "activation-readiness", "envelope": review["envelope"], "review_digest": review["review_digest"]})
				require.Equal(t, false, ready["readiness"].(map[string]any)["ready"])
				require.Empty(t, ready["readiness_digest"])
				require.Equal(t, before, sqlServerSuccessorLocalState(t, r.f.db))
				require.Equal(t, external, r.targetState(t), "readiness must not repair catalog damage or touch history/receipts")
			})
		}
	}
}
