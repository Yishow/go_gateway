package datalink

import (
	"context"
	"database/sql/driver"
	"net"
	"time"

	"github.com/lib/pq"
)

// lib/pq applies a context to dialing but can otherwise wait indefinitely for
// the startup handshake. Close the socket on cancellation only until the
// driver has completed that handshake; pooled connections must outlive it.
type postgresHandshakeConnector struct{ dsn string }

func (c postgresHandshakeConnector) Connect(ctx context.Context) (driver.Conn, error) {
	connector, err := pq.NewConnector(c.dsn)
	if err != nil {
		return nil, err
	}
	var stop func() bool
	dialer := &postgresHandshakeDialer{ctx: ctx, onConnect: func(conn net.Conn) {
		stop = context.AfterFunc(ctx, func() { _ = conn.Close() })
	}}
	connector.Dialer(dialer)
	conn, err := connector.Connect(ctx)
	if stop != nil {
		stop()
	}
	if ctx.Err() != nil {
		if err == nil {
			_ = conn.Close()
		}
		return nil, ctx.Err()
	}
	return conn, err
}

func (postgresHandshakeConnector) Driver() driver.Driver { return &pq.Driver{} }

type postgresHandshakeDialer struct {
	ctx       context.Context
	onConnect func(net.Conn)
}

func (d *postgresHandshakeDialer) Dial(network, address string) (net.Conn, error) {
	return d.DialContext(d.ctx, network, address)
}

func (d *postgresHandshakeDialer) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(d.ctx, timeout)
	defer cancel()
	return d.DialContext(ctx, network, address)
}

func (d *postgresHandshakeDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, network, address)
	if err == nil {
		d.onConnect(conn)
	}
	return conn, err
}
