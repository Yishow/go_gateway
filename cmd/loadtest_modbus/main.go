package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go-gateway/internal/datalink"
	"go-gateway/internal/datalink/collector"
	"go-gateway/internal/datalink/connector"
	_ "go-gateway/internal/datalink/connector/adapters"
	"go-gateway/internal/datalink/device"
	"go-gateway/internal/datalink/mapping"
	"go-gateway/internal/datalink/point"
	"go-gateway/internal/datalink/pollinggroup"
	runtime2 "go-gateway/internal/datalink/runtime"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/storage"
	"go-gateway/internal/datalink/tag"
	virtualmemory "go-gateway/internal/virtual/memory"
	virtualmodbus "go-gateway/internal/virtual/server/modbus"
)

type serverHandle struct {
	protocol schema.ProtocolType
	tcp      *virtualmodbus.Server
	udp      *virtualmodbus.UDPServer
	port     int
	expected uint16
}

func (s *serverHandle) stop() {
	if s.tcp != nil {
		_ = s.tcp.Stop()
	}
	if s.udp != nil {
		_ = s.udp.Stop()
	}
}

type report struct {
	StartedAt      time.Time          `json:"started_at"`
	FinishedAt     time.Time          `json:"finished_at"`
	Duration       string             `json:"duration"`
	TCPServers     int                `json:"tcp_servers"`
	UDPServers     int                `json:"udp_servers"`
	TotalServers   int                `json:"total_servers"`
	TotalRows      int64              `json:"total_rows"`
	GoodRows       int64              `json:"good_rows"`
	BadRows        int64              `json:"bad_rows"`
	MissingTagRows []string           `json:"missing_tag_rows"`
	MismatchCount  int64              `json:"mismatch_count"`
	MismatchByTag  map[string]int64   `json:"mismatch_by_tag"`
	ExpectedByTag  map[string]float64 `json:"expected_by_tag"`
	ExpectedByKey  map[string]float64 `json:"expected_by_key"`
	MismatchByKey  map[string]int64   `json:"mismatch_by_key"`
	RuntimeStats   runtime2.Stats     `json:"runtime_stats"`
	Pass           bool               `json:"pass"`
	FailReasons    []string           `json:"fail_reasons"`
	DatabaseFile   string             `json:"database_file"`
}

func main() {
	var (
		tcpCount     = flag.Int("tcp", 15, "TCP Modbus servers")
		udpCount     = flag.Int("udp", 15, "UDP Modbus servers")
		durationText = flag.String("duration", "10m", "test duration, e.g. 10m")
		intervalMs   = flag.Int("interval-ms", 500, "polling interval (ms)")
		dbPath       = flag.String("db", filepath.Join(os.TempDir(), "go_gateway_modbus_loadtest.db"), "sqlite db path")
		batchSize    = flag.Int("batch-size", 256, "batch writer size")
		flushMs      = flag.Int("flush-ms", 1000, "batch writer flush interval ms")
		reportPath   = flag.String("report", "", "optional json report output path")
		cleanDB      = flag.Bool("clean", true, "clean db file before run")
		valueBase    = flag.Int("value-base", 1000, "base register value for reproducible datasets")
		tagKeyPrefix = flag.String("tag-key-prefix", "bench/server", "tag key prefix for reproducible datasets")
	)
	flag.Parse()

	dur, err := time.ParseDuration(*durationText)
	if err != nil {
		log.Fatalf("invalid duration: %v", err)
	}
	if *tcpCount+*udpCount <= 0 {
		log.Fatalf("tcp + udp must be > 0")
	}
	if *valueBase < 0 || *valueBase > 65535 {
		log.Fatalf("value-base must be between 0 and 65535")
	}

	ctx := context.Background()
	if *cleanDB {
		_ = os.Remove(*dbPath)
	}

	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc", *dbPath)
	dbMgr := datalink.NewDBManager(datalink.DBConfig{
		Type:         datalink.DBTypeSQLite,
		DSN:          dsn,
		MaxOpenConns: 1,
		MaxIdleConns: 1,
	})
	if err := dbMgr.Connect(); err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer dbMgr.Close()
	db := dbMgr.DB()

	if err := datalink.NewMigrator().Migrate(db); err != nil {
		log.Fatalf("migrate failed: %v", err)
	}

	connMgr := connector.NewConnectionManager(connector.DefaultConnectionManagerConfig())
	defer connMgr.CloseAll()

	devRepo := device.NewSQLRepository(db)
	ptRepo := point.NewSQLRepository(db)
	tagRepo := tag.NewSQLRepository(db)
	mapRepo := mapping.NewSQLRepository(db)
	pgRepo := pollinggroup.NewSQLRepository(db)

	devSvc := device.NewService(devRepo, connMgr)
	ptSvc := point.NewService(ptRepo, pgRepo)
	tagSvc := tag.NewService(tagRepo)
	mapSvc := mapping.NewServiceWithTagResolver(mapRepo, tagSvc.GetByID)
	pgSvc := pollinggroup.NewService(pgRepo)

	schedulerCfg := collector.DefaultSchedulerConfig()
	schedulerCfg.ValueBufferSize = 20000
	scheduler := collector.NewScheduler(schedulerCfg, connMgr)

	sqliteWriter := storage.NewSQLiteWriter(db)
	bw := storage.NewBatchWriter(sqliteWriter, storage.BatchWriterConfig{
		BatchSize:     *batchSize,
		FlushInterval: time.Duration(*flushMs) * time.Millisecond,
		WriteTimeout:  5 * time.Second,
	})

	rt, err := runtime2.NewService(runtime2.Config{UpdatePointState: false}, runtime2.Dependencies{
		Scheduler:           scheduler,
		Writer:              bw,
		DeviceService:       devSvc,
		PointService:        ptSvc,
		MappingService:      mapSvc,
		TagService:          tagSvc,
		PollingGroupService: pgSvc,
	})
	if err != nil {
		log.Fatalf("create runtime failed: %v", err)
	}

	servers, err := startServers(*tcpCount, *udpCount, uint16(*valueBase))
	if err != nil {
		log.Fatalf("start servers failed: %v", err)
	}
	defer func() {
		for _, s := range servers {
			s.stop()
		}
	}()

	group, err := pgSvc.Create(ctx, pollinggroup.CreateRequest{
		Name:       "loadtest-30-modbus",
		IntervalMs: *intervalMs,
		Priority:   100,
	})
	if err != nil {
		log.Fatalf("create polling group failed: %v", err)
	}

	expectedByTag := make(map[string]float64)
	expectedByKey := make(map[string]float64)
	tagIDToKey := make(map[string]string)
	for i, srv := range servers {
		dev, err := devSvc.Create(ctx, device.CreateDeviceRequest{
			Name:        fmt.Sprintf("bench-dev-%02d", i),
			Protocol:    srv.protocol,
			Description: "local benchmark source",
			ConnectionConfig: map[string]interface{}{
				"host":     "127.0.0.1",
				"port":     srv.port,
				"slave_id": 1,
				"timeout":  2,
			},
		})
		if err != nil {
			log.Fatalf("create device %d failed: %v", i, err)
		}
		if err := devSvc.Activate(ctx, dev.ID); err != nil {
			log.Fatalf("activate device %s failed: %v", dev.ID, err)
		}

		gid := group.ID
		pt, err := ptSvc.Create(ctx, point.CreatePointRequest{
			DeviceID:       dev.ID,
			Name:           fmt.Sprintf("bench-pt-%02d", i),
			Address:        "40001",
			Function:       "03",
			DataType:       schema.DataTypeUint16,
			Mode:           schema.PointModeReadOnly,
			PollingGroupID: &gid,
		})
		if err != nil {
			log.Fatalf("create point %d failed: %v", i, err)
		}

		tagKey := fmt.Sprintf("%s/%02d", *tagKeyPrefix, i)
		tg, err := tagSvc.Create(ctx, tag.CreateTagRequest{
			Key:         tagKey,
			DisplayName: fmt.Sprintf("Bench Server %02d", i),
			DataType:    schema.DataTypeUint16,
			Unit:        "count",
		})
		if err != nil {
			log.Fatalf("create tag %d failed: %v", i, err)
		}
		if err := tagSvc.Activate(ctx, tg.ID); err != nil {
			log.Fatalf("activate tag %s failed: %v", tg.ID, err)
		}

		if _, err := mapSvc.Create(ctx, mapping.CreateMappingRequest{
			PointID:           pt.ID,
			TagID:             tg.ID,
			TransformPipeline: []schema.TransformStep{},
		}); err != nil {
			log.Fatalf("create mapping %d failed: %v", i, err)
		}

		expectedByTag[tg.ID] = float64(srv.expected)
		expectedByKey[tagKey] = float64(srv.expected)
		tagIDToKey[tg.ID] = tagKey
	}

	if err := rt.RefreshMappings(ctx); err != nil {
		log.Fatalf("refresh mappings failed: %v", err)
	}

	startedAt := time.Now()
	if err := rt.Start(ctx); err != nil {
		log.Fatalf("runtime start failed: %v", err)
	}

	log.Printf("loadtest started: tcp=%d udp=%d duration=%s interval=%dms", *tcpCount, *udpCount, dur, *intervalMs)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	deadline := time.NewTimer(dur)
	defer deadline.Stop()

	running := true
	for running {
		select {
		case <-ticker.C:
			st := rt.Snapshot()
			log.Printf("progress collected=%d written=%d write_err=%d mapping_err=%d", st.CollectedTotal, st.WriteSuccess, st.WriteError, st.MappingError)
		case <-deadline.C:
			running = false
		}
	}

	if err := rt.Stop(ctx); err != nil {
		log.Fatalf("runtime stop failed: %v", err)
	}
	finishedAt := time.Now()

	rep, err := verify(ctx, db, expectedByTag, expectedByKey, tagIDToKey, rt.Snapshot(), startedAt, finishedAt, *tcpCount, *udpCount, *dbPath)
	if err != nil {
		log.Fatalf("verify failed: %v", err)
	}

	printReport(rep)

	if *reportPath != "" {
		if err := writeReport(*reportPath, rep); err != nil {
			log.Fatalf("write report failed: %v", err)
		}
		log.Printf("report written: %s", *reportPath)
	}

	if !rep.Pass {
		os.Exit(1)
	}
}

func startServers(tcpCount, udpCount int, valueBase uint16) ([]*serverHandle, error) {
	servers := make([]*serverHandle, 0, tcpCount+udpCount)

	for i := 0; i < tcpCount; i++ {
		bank := virtualmemory.NewMemoryBank(2048)
		expected := valueBase + uint16(i)
		if err := bank.WriteWord(0, expected); err != nil {
			return nil, err
		}

		srv := virtualmodbus.NewServer(bank)
		if err := srv.Start(0); err != nil {
			return nil, err
		}
		servers = append(servers, &serverHandle{protocol: schema.ProtocolModbusTCP, tcp: srv, port: srv.Port(), expected: expected})
	}

	for i := 0; i < udpCount; i++ {
		bank := virtualmemory.NewMemoryBank(2048)
		expected := valueBase + uint16(tcpCount+i)
		if err := bank.WriteWord(0, expected); err != nil {
			return nil, err
		}

		srv := virtualmodbus.NewUDPServer(bank)
		if err := srv.Start(0); err != nil {
			return nil, err
		}
		servers = append(servers, &serverHandle{protocol: schema.ProtocolModbusUDP, udp: srv, port: srv.Port(), expected: expected})
	}

	return servers, nil
}

func verify(ctx context.Context, db *sql.DB, expectedByTag map[string]float64, expectedByKey map[string]float64, tagIDToKey map[string]string, stats runtime2.Stats, startedAt, finishedAt time.Time, tcpCount, udpCount int, dbPath string) (*report, error) {
	rep := &report{
		StartedAt:     startedAt,
		FinishedAt:    finishedAt,
		Duration:      finishedAt.Sub(startedAt).String(),
		TCPServers:    tcpCount,
		UDPServers:    udpCount,
		TotalServers:  tcpCount + udpCount,
		MismatchByTag: make(map[string]int64),
		ExpectedByTag: expectedByTag,
		ExpectedByKey: expectedByKey,
		MismatchByKey: make(map[string]int64),
		RuntimeStats:  stats,
		DatabaseFile:  dbPath,
	}

	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM timeseries`).Scan(&rep.TotalRows); err != nil {
		return nil, err
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM timeseries WHERE quality = 'good'`).Scan(&rep.GoodRows); err != nil {
		return nil, err
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM timeseries WHERE quality != 'good'`).Scan(&rep.BadRows); err != nil {
		return nil, err
	}

	for tagID, expectedVal := range expectedByTag {
		var cnt int64
		if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM timeseries WHERE tag_id = ?`, tagID).Scan(&cnt); err != nil {
			return nil, err
		}
		if cnt == 0 {
			rep.MissingTagRows = append(rep.MissingTagRows, tagID)
		}

		var mismatch int64
		if err := db.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM timeseries
			WHERE tag_id = ?
			  AND quality = 'good'
			  AND (value_num IS NULL OR ABS(value_num - ?) > ?)
		`, tagID, expectedVal, 0.0001).Scan(&mismatch); err != nil {
			return nil, err
		}
		rep.MismatchByTag[tagID] = mismatch
		if key, ok := tagIDToKey[tagID]; ok {
			rep.MismatchByKey[key] = mismatch
		}
		rep.MismatchCount += mismatch
	}

	reasons := make([]string, 0)
	if rep.TotalRows == 0 {
		reasons = append(reasons, "timeseries 無資料")
	}
	if len(rep.MissingTagRows) > 0 {
		reasons = append(reasons, fmt.Sprintf("%d 個 tag 無資料", len(rep.MissingTagRows)))
	}
	if rep.MismatchCount > 0 {
		reasons = append(reasons, fmt.Sprintf("資料不一致 %d 筆", rep.MismatchCount))
	}
	if stats.WriteError > 0 {
		reasons = append(reasons, fmt.Sprintf("write error=%d", stats.WriteError))
	}
	if stats.MappingError > 0 {
		reasons = append(reasons, fmt.Sprintf("mapping error=%d", stats.MappingError))
	}
	if stats.CollectedTotal == 0 || stats.WriteSuccess == 0 {
		reasons = append(reasons, "未收到有效採集/寫入")
	}

	rep.FailReasons = reasons
	rep.Pass = len(reasons) == 0
	return rep, nil
}

func printReport(rep *report) {
	fmt.Println("================ LOAD TEST REPORT ================")
	fmt.Printf("Duration: %s\n", rep.Duration)
	fmt.Printf("Servers : TCP=%d UDP=%d Total=%d\n", rep.TCPServers, rep.UDPServers, rep.TotalServers)
	fmt.Printf("Rows    : total=%d good=%d bad=%d\n", rep.TotalRows, rep.GoodRows, rep.BadRows)
	fmt.Printf("Runtime : collected=%d write_ok=%d write_err=%d map_err=%d\n",
		rep.RuntimeStats.CollectedTotal,
		rep.RuntimeStats.WriteSuccess,
		rep.RuntimeStats.WriteError,
		rep.RuntimeStats.MappingError,
	)
	fmt.Printf("Mismatch: %d\n", rep.MismatchCount)
	if len(rep.MissingTagRows) > 0 {
		fmt.Printf("Missing : %d tags\n", len(rep.MissingTagRows))
	}
	if rep.Pass {
		fmt.Println("Result  : PASS ✅")
	} else {
		fmt.Println("Result  : FAIL ❌")
		for _, r := range rep.FailReasons {
			fmt.Printf("  - %s\n", r)
		}
	}
	fmt.Println("==================================================")
}

func writeReport(path string, rep *report) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(rep, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}
