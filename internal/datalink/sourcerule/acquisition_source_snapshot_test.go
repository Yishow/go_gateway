package sourcerule

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAcquisitionSourceSQLSnapshotEnumeratesStandaloneAndLayout(t *testing.T) {
	db, _ := shareRecoveryFixture(t)
	_, err := db.ExecContext(t.Context(), `INSERT INTO points(id,device_id,name,address,data_type,enabled) VALUES('standalone','d','standalone','40005','uint16',0)`)
	require.NoError(t, err)
	_, err = db.ExecContext(t.Context(), `UPDATE source_rules SET layout_mode='mixed',mixed_items='[{"address":"40001"}]'`)
	require.NoError(t, err)
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	rows, err := ReadAcquisitionSourceRowsTx(t.Context(), tx, "d")
	require.NoError(t, err)
	require.Len(t, rows.Points, 2)
	require.Equal(t, "mixed", rows.Rules[0].LayoutMode)
	require.NotEmpty(t, rows.Rules[0].MixedItems)
	_, reason := AcquisitionSourceExclusions(rows)
	require.Equal(t, "source-layout-unsupported", reason)
}

func TestAcquisitionSourceSQLSnapshotRejectsForeignOwner(t *testing.T) {
	db, _ := shareRecoveryFixture(t)
	for _, query := range []string{
		`INSERT INTO devices(id,name,protocol,status) VALUES('f','f','modbus_tcp','active')`,
		`INSERT INTO source_rules(id,device_id,start_address,count,data_type,naming_prefix) VALUES('fr','f','40001',1,'uint16','f')`,
		`UPDATE source_rule_links SET rule_id='fr' WHERE id='l'`,
	} {
		_, err := db.ExecContext(t.Context(), query)
		require.NoError(t, err)
	}
	tx, err := db.BeginTx(t.Context(), nil)
	require.NoError(t, err)
	defer tx.Rollback()
	_, err = ReadAcquisitionSourceRowsTx(t.Context(), tx, "d")
	require.ErrorIs(t, err, ErrAcquisitionSourceSnapshot)
}
