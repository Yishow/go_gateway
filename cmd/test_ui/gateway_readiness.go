package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"go-gateway/internal/apphost"
	"go-gateway/internal/desktop"
)

// probeOwnHTTP confirms the reserved listener serves this artifact, before a
// ready event or browser action. It never uses ambient HTTP proxy settings.
func probeOwnHTTP(ctx context.Context, address string, expected []byte) error {
	origin, _, err := apphost.ListenerURL(address)
	if err != nil {
		return err
	}
	dialer := net.Dialer{Timeout: 2 * time.Second}
	transport := &http.Transport{Proxy: nil, DialContext: dialer.DialContext, DisableKeepAlives: true, ResponseHeaderTimeout: 2 * time.Second}
	defer transport.CloseIdleConnections()
	client := http.Client{Transport: transport, Timeout: 3 * time.Second}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, origin+"/studio/v2", http.NoBody)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, int64(len(expected))+1))
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK || !bytes.Equal(body, expected) {
		return errors.New("embedded service readiness response mismatch")
	}
	return nil
}

// verifyOpenedDatabase is a second identity check before any migration. Windows
// additionally retains ancestor and file ownership handles across this reopen.
func verifyOpenedDatabase(ctx context.Context, db *sql.DB, expected desktop.Identity) error {
	var sequence int
	var schemaName, actualPath string
	if err := db.QueryRowContext(ctx, "PRAGMA database_list").Scan(&sequence, &schemaName, &actualPath); err != nil {
		return err
	}
	identity, err := desktop.CanonicalDatabase(actualPath)
	if err != nil {
		return err
	}
	if identity.ID != expected.ID {
		return errors.New("opened database identity differs from owner")
	}
	return nil
}
