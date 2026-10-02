package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"go-gateway/internal/datalink/dbtarget"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

// tcpProxy forwards to a real PostgreSQL so a test can take the destination
// offline (refusing and dropping connections) and bring it back on the same port.
type tcpProxy struct {
	t        *testing.T
	target   string
	addr     string
	mu       sync.Mutex
	listener net.Listener
	conns    []net.Conn
}

func newTCPProxy(t *testing.T, target string) *tcpProxy {
	t.Helper()
	p := &tcpProxy{t: t, target: target}
	p.up("127.0.0.1:0")
	t.Cleanup(p.down)
	return p
}

func (p *tcpProxy) up(addr string) {
	p.t.Helper()
	var listener net.Listener
	var err error
	for attempt := 0; attempt < 50; attempt++ { // the port may linger briefly after down()
		listener, err = net.Listen("tcp", addr) //nolint:noctx // test proxy: the listener lifetime is managed by up/down
		if err == nil {
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	require.NoError(p.t, err)
	p.mu.Lock()
	p.listener, p.addr = listener, listener.Addr().String()
	p.mu.Unlock()
	go func() {
		for {
			client, err := listener.Accept()
			if err != nil {
				return
			}
			upstream, err := net.Dial("tcp", p.target) //nolint:noctx // test proxy forwarding to the disposable PostgreSQL
			if err != nil {
				_ = client.Close()
				continue
			}
			p.mu.Lock()
			p.conns = append(p.conns, client, upstream)
			p.mu.Unlock()
			go func() { _, _ = io.Copy(upstream, client); _ = upstream.Close() }()
			go func() { _, _ = io.Copy(client, upstream); _ = client.Close() }()
		}
	}()
}

func (p *tcpProxy) down() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.listener != nil {
		_ = p.listener.Close()
		p.listener = nil
	}
	for _, conn := range p.conns {
		_ = conn.Close()
	}
	p.conns = nil
}

func (p *tcpProxy) host() (host, port string) {
	h, pt, err := net.SplitHostPort(p.addr)
	require.NoError(p.t, err)
	return h, pt
}

func postgresEnvFields(t *testing.T) (dsn string, fields map[string]string) {
	t.Helper()
	dsn = strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("Skipping PostgreSQL group delivery test: POSTGRES_DSN not set")
	}
	fields = map[string]string{}
	for _, part := range strings.Fields(dsn) {
		if key, value, ok := strings.Cut(part, "="); ok {
			fields[key] = value
		}
	}
	return dsn, fields
}

func pgCount(t *testing.T, admin *sql.DB, pgSchema, table string) int {
	t.Helper()
	var n int
	require.NoError(t, admin.QueryRowContext(t.Context(), `SELECT COUNT(*) FROM "`+pgSchema+`".`+table).Scan(&n))
	return n
}

func TestProductionGroupOutageRecoveryPostgresDestinationReceiptsAndRestart(t *testing.T) {
	dsn, fields := postgresEnvFields(t)
	admin, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = admin.Close() })
	require.NoError(t, admin.PingContext(t.Context()))

	schemas := map[string]string{}
	proxies := map[string]*tcpProxy{}
	setup := func(t *testing.T, env *outageEnv, name string) (kind, configJSON, tableSchema string) {
		t.Helper()
		pgSchema := fmt.Sprintf("gw_pipe_%s_%d", strings.ToLower(name), time.Now().UnixNano())
		_, err := admin.ExecContext(t.Context(), `CREATE SCHEMA "`+pgSchema+`"`)
		require.NoError(t, err)
		t.Cleanup(func() { _, _ = admin.ExecContext(context.Background(), `DROP SCHEMA "`+pgSchema+`" CASCADE`) })
		_, err = admin.ExecContext(t.Context(), `CREATE TABLE "`+pgSchema+`".readings (temperature DOUBLE PRECISION NOT NULL, pressure BIGINT NOT NULL)`)
		require.NoError(t, err)
		require.NoError(t, dbtarget.CreateEffectReceiptTable(t.Context(), admin, schema.DatabaseConnectorKindPostgres, pgSchema))
		port := fields["port"]
		if port == "" {
			port = "5432"
		}
		proxy := newTCPProxy(t, net.JoinHostPort(fields["host"], port))
		host, proxyPort := proxy.host()
		schemas[name], proxies[name] = pgSchema, proxy
		config, err := json.Marshal(map[string]string{
			"host": host, "port": proxyPort, "user": fields["user"], "password": fields["password"],
			"database": fields["dbname"], "sslmode": "disable",
		})
		require.NoError(t, err)
		return "postgres", string(config), pgSchema
	}
	env := newOutageEnvWith(t, setup)
	env.dedupe = "receipt"
	groupA := env.createAndApply(t, "A")
	groupB := env.createAndApply(t, "B")
	ctx := t.Context()
	groups := []*workspace.WriteGroup{groupA, groupB}

	clock := &testClock{}
	resolved, err := env.services.writeGroups.ResolveAppliedAt(ctx, groupA.ID, time.Now().UTC().Add(time.Minute))
	require.NoError(t, err)
	base := resolved.EffectiveAt.Truncate(10 * time.Second).Add(10 * time.Second)
	clock.set(base.Add(time.Second))
	first := env.pipeline(clock, "node-1/first")
	require.NoError(t, first.Start(ctx))
	for _, status := range first.Status() {
		require.Equal(t, "active", status.State, "group %s blocked: %s", status.GroupID, status.Reason)
	}

	env.feed(t, first, clock, groups, base, 21.5, 100)
	require.Eventually(t, func() bool {
		return pgCount(t, admin, schemas["A"], "readings") == 1 && pgCount(t, admin, schemas["B"], "readings") == 1
	}, 15*time.Second, 25*time.Millisecond)

	// Destination A goes offline: connections are refused and dropped.
	proxies["A"].down()
	env.feed(t, first, clock, groups, base.Add(10*time.Second), 22.5, 200)
	env.feed(t, first, clock, groups, base.Add(20*time.Second), 23.5, 300)
	require.Eventually(t, func() bool { return pgCount(t, admin, schemas["B"], "readings") == 3 }, 15*time.Second, 25*time.Millisecond,
		"the healthy PostgreSQL destination keeps receiving while the other is down")
	require.Equal(t, 1, pgCount(t, admin, schemas["A"], "readings"), "nothing reaches the offline destination")
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state != 'sql_committed'") == 2 }, 10*time.Second, 25*time.Millisecond)
	view, err := first.Delivery(ctx, groupA.ID)
	require.NoError(t, err)
	require.Equal(t, 1, view.Stages.SQLCommitted)
	require.Equal(t, 2, view.Stages.Queued+view.Stages.Retrying, "the outage rows are accepted and waiting, not written")

	// The gateway restarts during the outage.
	require.NoError(t, first.Stop(2*time.Second))
	second := env.pipeline(clock, "node-1/second")
	require.NoError(t, second.Start(ctx))
	require.Equal(t, 1, pgCount(t, admin, schemas["A"], "readings"))

	// A comes back: the backlog is delivered once, in order, with exact values.
	proxies["A"].up(proxies["A"].addr)
	require.Eventually(t, func() bool { return pgCount(t, admin, schemas["A"], "readings") == 3 }, 20*time.Second, 25*time.Millisecond)
	require.Eventually(t, func() bool { return outboxCount(t, env.db, "connector-A", "state = 'sql_committed'") == 3 }, 10*time.Second, 25*time.Millisecond)
	require.NoError(t, second.Stop(2*time.Second))
	for _, name := range []string{"A", "B"} {
		rows, err := admin.QueryContext(ctx, `SELECT temperature, pressure FROM "`+schemas[name]+`".readings ORDER BY temperature`)
		require.NoError(t, err)
		var got [][2]float64
		for rows.Next() {
			var temperature float64
			var pressure int64
			require.NoError(t, rows.Scan(&temperature, &pressure))
			got = append(got, [2]float64{temperature, float64(pressure)})
		}
		require.NoError(t, rows.Close())
		require.Equal(t, [][2]float64{{21.5, 100}, {22.5, 200}, {23.5, 300}}, got, "destination %s holds each row exactly once", name)
		require.Equal(t, 3, pgCount(t, admin, schemas[name], dbtarget.EffectReceiptTable), "each row has its destination-side receipt")
	}
}
