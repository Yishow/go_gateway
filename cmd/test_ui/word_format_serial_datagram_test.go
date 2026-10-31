package main

import (
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
)

func wordDatagramPeer(t *testing.T) (*wordWirePLC, string) {
	t.Helper()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	require.NoError(t, err)
	plc := &wordWirePLC{words: map[int]uint16{}}
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			raw := make([]byte, 260)
			n, addr, err := conn.ReadFromUDP(raw)
			if err != nil {
				return
			}
			if n != 12 {
				continue
			}
			start, count := int(binary.BigEndian.Uint16(raw[8:10])), int(binary.BigEndian.Uint16(raw[10:12]))
			words, bad := plc.read(start, count)
			response := append([]byte(nil), raw[:7]...)
			response = append(response, raw[7], byte(count*2))
			for _, w := range words {
				response = append(response, byte(w>>8), byte(w))
			}
			if bad {
				response = response[:len(response)-1]
			}
			binary.BigEndian.PutUint16(response[4:6], uint16(len(response)-6))
			_, _ = conn.WriteToUDP(response, addr)
		}
	}()
	t.Cleanup(func() { require.NoError(t, conn.Close()); <-done })
	return plc, fmt.Sprintf(`{"host":"127.0.0.1","port":%d,"slave_id":1,"timeout":1}`, conn.LocalAddr().(*net.UDPAddr).Port)
}

// Only the byte transport peer is simulated; factory, manager, probe, source
// HTTP route, persistent point, scheduler, adapter/client and ingestor are real.
func newWordTransportFixture(t *testing.T, protocol schema.ProtocolType, plc *wordWirePLC, config string) wordProductionFixture {
	t.Helper()
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "word-transport.db"))
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	t.Cleanup(func() { require.NoError(t, db.Close()) })
	require.NoError(t, datalink.NewMigrator().Migrate(db))
	_, err = db.Exec(`INSERT INTO devices(id,name,description,protocol,status,connection_config,readiness_status) VALUES('word-device','PLC','',?,'draft',?,'{}')`, protocol, config)
	require.NoError(t, err)
	manager := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	t.Cleanup(func() { require.NoError(t, manager.CloseAll()) })
	services, err := wireGatewayServices(db, manager)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, services.runtime.Stop(context.Background())) })
	require.NoError(t, services.device.ProbeAndActivate(t.Context(), "word-device"))
	_, err = services.workspace.AttachDevice(t.Context(), "word-device")
	require.NoError(t, err)
	router := wordProductionRouter(services, db)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/devices/word-device/activate", nil))
	require.Equal(t, 200, response.Code, response.Body.String())
	return wordProductionFixture{services, db, plc, manager, router}
}

func exerciseWordTransportMatrix(t *testing.T, f wordProductionFixture, protocol schema.ProtocolType) {
	t.Helper()
	type expected struct {
		point  *schema.Point
		value  any
		format string
	}
	var cases []expected
	for _, dtype := range []schema.DataType{schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32, schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64} {
		for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
			// createRule's protocol switch only selects the address spelling. All
			// configuration/authoring/resolution uses the persisted RTU/UDP device.
			_, p, want := f.createRule(t, schema.ProtocolModbusTCP, len(cases), dtype, format)
			cases = append(cases, expected{p, want, format})
		}
	}
	require.NoError(t, f.services.runtime.Start(t.Context()))
	for _, tc := range cases {
		t.Run(string(tc.point.DataType)+"/"+tc.format, func(t *testing.T) {
			raw, err := json.Marshal(tc.value)
			require.NoError(t, err)
			require.Eventually(t, func() bool {
				p, e := f.services.point.GetByID(t.Context(), tc.point.ID)
				return e == nil && p.LastValue != nil && *p.LastValue == string(raw) && p.LastError == ""
			}, 10*time.Second, 20*time.Millisecond, "scheduled actual adapter/ingestor")
			values := f.services.scheduler.PollNowContext(t.Context(), []string{tc.point.ID})
			require.Len(t, values, 1)
			cv := values[0]
			require.Equal(t, schema.QualityGood, cv.Quality)
			require.Equal(t, tc.value, cv.Value)
			require.Equal(t, tc.point.RevisionID, cv.PointRevision)
			require.NotZero(t, cv.InstallationLease)
			require.NotEmpty(t, cv.InstallationScope)
			profile, ok := f.services.scheduler.CapturePointReadProfile(tc.point.ID, cv.InstallationLease)
			require.True(t, ok)
			require.Equal(t, protocol, profile.IdentityView().Protocol)
			require.Equal(t, tc.format, profile.Resolution().EffectiveFormat)
			d, err := f.services.device.GetByID(t.Context(), tc.point.DeviceID)
			require.NoError(t, err)
			fingerprint, err := measurement.AcquisitionConfigFingerprint(*d, *tc.point, profile)
			require.NoError(t, err)
			require.Equal(t, fingerprint, cv.ConfigFingerprint)
		})
	}
}

func TestProductionWordFormatUDPToIngestorMatrix(t *testing.T) {
	plc, config := wordDatagramPeer(t)
	exerciseWordTransportMatrix(t, newWordTransportFixture(t, schema.ProtocolModbusUDP, plc, config), schema.ProtocolModbusUDP)
}
