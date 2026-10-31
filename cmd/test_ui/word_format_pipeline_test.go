package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go-gateway/internal/api"
	appconfig "go-gateway/internal/config"
	"go-gateway/internal/datalink/modbusshare"
	"math"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/measurement"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/sourcerule"
)

type wordProductionFixture struct {
	services gatewayServices
	db       *sql.DB
	plc      *wordWirePLC
	manager  *connector.ConnectionManager
	router   http.Handler
}

func newWordProductionFixture(t *testing.T, protocol schema.ProtocolType) wordProductionFixture {
	t.Helper()
	plc, config := wordWireServer(t, protocol)
	db, err := sql.Open("sqlite", filepath.Join(t.TempDir(), "word.db"))
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
	request := httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/devices/word-device/activate", nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	require.Equal(t, 200, response.Code, response.Body.String())
	return wordProductionFixture{services, db, plc, manager, router}
}
func wordCanonicalValue(dtype schema.DataType) (uint64, any, int) {
	switch dtype {
	case schema.DataTypeInt32:
		return 0x11223344, int32(0x11223344), 2
	case schema.DataTypeUint32:
		return 0x11223344, uint32(0x11223344), 2
	case schema.DataTypeFloat32:
		return uint64(math.Float32bits(12.5)), float32(12.5), 2
	case schema.DataTypeInt64:
		return 0x112233445566, int64(0x112233445566), 4
	case schema.DataTypeUint64:
		return 0x112233445566, uint64(0x112233445566), 4
	default:
		return math.Float64bits(22.25), float64(22.25), 4
	}
}

// Independent wire encoding: byte swap within each word and full word reversal.
func wordWireEncoding(bits uint64, count int, format string) []uint16 {
	raw := make([]byte, 8)
	binary.BigEndian.PutUint64(raw, bits)
	raw = raw[8-count*2:]
	words := make([]uint16, count)
	for i := range words {
		words[i] = binary.BigEndian.Uint16(raw[2*i:])
	}
	if format == "BADC" || format == "DCBA" {
		for i, w := range words {
			words[i] = w<<8 | w>>8
		}
	}
	if format == "CDAB" || format == "DCBA" {
		for i, j := 0, len(words)-1; i < j; i, j = i+1, j-1 {
			words[i], words[j] = words[j], words[i]
		}
	}
	return words
}
func (f wordProductionFixture) createRule(t *testing.T, protocol schema.ProtocolType, index int, dtype schema.DataType, format string, addressPrefix ...string) (*schema.SourceRule, *schema.Point, any) {
	t.Helper()
	start := index * 8
	address := fmt.Sprintf("D%d", start)
	if len(addressPrefix) > 0 {
		address = fmt.Sprintf("%s%d", addressPrefix[0], start)
	}
	if protocol == schema.ProtocolModbusTCP {
		address = fmt.Sprintf("%d", 40001+start)
	}
	bits, want, count := wordCanonicalValue(dtype)
	f.plc.set(start, wordWireEncoding(bits, count, format))
	req := sourcerule.CreateRuleRequest{DeviceID: "word-device", StartAddress: address, Count: 1, DataType: dtype, NamingPrefix: fmt.Sprintf("word%d", index), Enabled: true, DataFormat: format}
	var rule *schema.SourceRule
	var err error
	if !f.services.modbusShare.Settings().Enabled {
		payload, e := json.Marshal(req)
		require.NoError(t, e)
		request := httptest.NewRequest("POST", appconfig.Get().API.BasePath+"/datalink/source-rules", bytes.NewReader(payload))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		f.router.ServeHTTP(response, request)
		require.Equal(t, http.StatusCreated, response.Code, response.Body.String())
		var body struct {
			Data struct {
				ID string `json:"id"`
			} `json:"data"`
		}
		require.NoError(t, json.Unmarshal(response.Body.Bytes(), &body))
		rule, err = f.services.sourceRule.GetByID(t.Context(), body.Data.ID)
	} else {
		rule, err = f.services.sourceRule.Create(t.Context(), req)
	}
	require.NoError(t, err)
	links, err := f.services.sourceRule.ListLinks(t.Context(), rule.ID)
	require.NoError(t, err)
	require.Len(t, links, 1)
	p, err := f.services.point.GetByID(t.Context(), links[0].PointID)
	require.NoError(t, err)
	require.Equal(t, schema.PointFormatInherit, p.FormatMode)
	require.True(t, p.Enabled, "source requires actual successful activation probe")
	return rule, p, want
}
func TestProductionWordFormatRuleToIngestorMatrix(t *testing.T) {
	for _, protocol := range []schema.ProtocolType{schema.ProtocolModbusTCP, schema.ProtocolMC3E, schema.ProtocolFatekFBs} {
		t.Run(string(protocol), func(t *testing.T) {
			f := newWordProductionFixture(t, protocol)
			type expected struct {
				p      *schema.Point
				want   any
				format string
			}
			var cases []expected
			prefixes := []string{"D"}
			for _, prefix := range prefixes {
				for _, dtype := range []schema.DataType{schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32, schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64} {
					for _, format := range []string{"ABCD", "BADC", "CDAB", "DCBA"} {
						_, p, want := f.createRule(t, protocol, len(cases), dtype, format, prefix)
						cases = append(cases, expected{p, want, format})
					}
				}
			}

			require.NoError(t, f.services.runtime.Start(t.Context()))
			for _, tc := range cases {
				raw, err := json.Marshal(tc.want)
				require.NoError(t, err)
				var lastState string
				t.Cleanup(func() {
					if t.Failed() {
						t.Log(lastState)
					}
				})
				require.Eventually(t, func() bool {
					p, e := f.services.point.GetByID(t.Context(), tc.p.ID)
					if e == nil {
						lastState = fmt.Sprintf("value=%v error=%s read=%v", p.LastValue, p.LastError, p.LastReadAt)
						if p.LastValue != nil {
							lastState += " raw=" + *p.LastValue
						}
					}
					return e == nil && p.LastValue != nil && *p.LastValue == string(raw) && p.LastError == ""
				}, 10*time.Second, 20*time.Millisecond, "scheduled actual ingestor %s/%s", tc.p.DataType, tc.format)
				values := f.services.scheduler.PollNowContext(t.Context(), []string{tc.p.ID})
				require.Len(t, values, 1)
				cv := values[0]
				require.Equal(t, schema.QualityGood, cv.Quality)
				require.Equal(t, tc.want, cv.Value)
				require.Equal(t, tc.p.RevisionID, cv.PointRevision)
				require.NotZero(t, cv.InstallationLease)
				require.NotEmpty(t, cv.InstallationScope)
				profile, ok := f.services.scheduler.CapturePointReadProfile(tc.p.ID, cv.InstallationLease)
				require.True(t, ok)
				require.Equal(t, tc.format, profile.Resolution().EffectiveFormat)
				d, err := f.services.device.GetByID(t.Context(), tc.p.DeviceID)
				require.NoError(t, err)
				fingerprint, err := measurement.AcquisitionConfigFingerprint(*d, *tc.p, profile)
				require.NoError(t, err)
				require.Equal(t, fingerprint, cv.ConfigFingerprint)
			}
		})
	}
}

func wordProductionRouter(s gatewayServices, db *sql.DB) http.Handler {
	revisions := modbusshare.NewSQLWorkspaceRevisionStore(db)
	reconciler := modbusshare.NewReconciler(s.modbusShare, revisions)
	configureShareRuntimeReconciler(s.sourceRule, s.workspace, s.modbusShare, revisions, reconciler)
	return api.NewRouter(&api.DatalinkServices{Device: s.device, Point: s.point, Tag: s.tag, Mapping: s.mapping, PollingGroup: s.pollingGroup, Settings: s.settings, ModbusShare: s.modbusShare, ModbusShareReconciler: reconciler, Scheduler: s.scheduler, Runtime: s.runtime, DBTarget: s.dbTarget, DBMapping: s.dbMapping, SourceRule: s.sourceRule, RecordingPlan: s.recordingPlan, Workspace: s.workspace, WriteGroups: s.writeGroups, WriteGroupDelivery: s.groupPipe, WriteGroupTestWrite: s.groupTestWrite, Audit: s.audit})
}
