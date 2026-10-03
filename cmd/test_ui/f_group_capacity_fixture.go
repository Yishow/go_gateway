//go:build f_write_group_fixture

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"go-gateway/internal/datalink/groupdelivery"
)

const (
	fixtureCapacityKindDiskFull = "disk_full"
	fixtureGroupQuotaEnv        = "F_FIXTURE_GROUP_QUOTA_BYTES"
	fixtureGlobalQuotaEnv       = "F_FIXTURE_GLOBAL_QUOTA_BYTES"
	fixtureQuotaMaxBytes        = int64(1 << 40)
)

type fixtureCapacityRequest struct {
	Kind    string `json:"kind"`
	Enabled bool   `json:"enabled"`
}

type fixtureCapacityState struct {
	Kind         string `json:"kind"`
	Enabled      bool   `json:"enabled"`
	PageCount    int64  `json:"page_count"`
	MaxPageCount int64  `json:"max_page_count"`
}

type fixtureCapacityController struct {
	db                   *sql.DB
	mu                   sync.Mutex
	enabled              bool
	originalMaxPageCount int64
}

func parseFixtureQuotaBytes(raw, envName string, fallback int64) (int64, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || value <= 0 || value > fixtureQuotaMaxBytes {
		return 0, fmt.Errorf("%s must be a positive decimal byte count within the fixture limit", envName)
	}
	return value, nil
}

func fixtureQuotaConfig() (groupdelivery.QuotaConfig, error) {
	global, err := parseFixtureQuotaBytes(os.Getenv(fixtureGlobalQuotaEnv), fixtureGlobalQuotaEnv, groupDeliveryQuotaBytes)
	if err != nil {
		return groupdelivery.QuotaConfig{}, err
	}
	group, err := parseFixtureQuotaBytes(os.Getenv(fixtureGroupQuotaEnv), fixtureGroupQuotaEnv, 0)
	if err != nil {
		return groupdelivery.QuotaConfig{}, err
	}
	return groupdelivery.QuotaConfig{GlobalMaxBytes: global, GroupMaxBytes: group}, nil
}

func preflightFixtureCapacity() error {
	_, err := fixtureQuotaConfig()
	if err != nil {
		return err
	}
	return preflightFixtureRetries()
}

func groupDeliveryQuotaConfig() groupdelivery.QuotaConfig {
	config, err := fixtureQuotaConfig()
	if err != nil {
		return groupdelivery.QuotaConfig{}
	}
	return config
}

func newFixtureCapacityController(db *sql.DB) *fixtureCapacityController {
	return &fixtureCapacityController{db: db}
}

func (c *fixtureCapacityController) configure(ctx context.Context, enabled bool) (fixtureCapacityState, error) {
	if c == nil || c.db == nil {
		return fixtureCapacityState{}, errors.New("fixture capacity controller is unavailable")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	if enabled && !c.enabled {
		pageCount, maxPageCount, err := c.pageCounts(ctx)
		if err != nil {
			return fixtureCapacityState{}, err
		}
		if maxPageCount <= 0 {
			return fixtureCapacityState{}, errors.New("fixture database page limit is invalid")
		}
		if err := c.setMaxPageCount(ctx, pageCount); err != nil {
			return fixtureCapacityState{}, err
		}
		c.originalMaxPageCount = maxPageCount
		c.enabled = true
	}
	if !enabled && c.enabled {
		if err := c.setMaxPageCount(ctx, c.originalMaxPageCount); err != nil {
			return fixtureCapacityState{}, err
		}
		c.originalMaxPageCount = 0
		c.enabled = false
	}
	pageCount, maxPageCount, err := c.pageCounts(ctx)
	if err != nil {
		return fixtureCapacityState{}, err
	}
	return fixtureCapacityState{
		Kind: fixtureCapacityKindDiskFull, Enabled: c.enabled,
		PageCount: pageCount, MaxPageCount: maxPageCount,
	}, nil
}

func (c *fixtureCapacityController) close(ctx context.Context) error {
	if c == nil {
		return nil
	}
	_, err := c.configure(ctx, false)
	return err
}

func (c *fixtureCapacityController) pageCounts(ctx context.Context) (pageCount, maxPageCount int64, err error) {
	if err := c.db.QueryRowContext(ctx, "PRAGMA page_count").Scan(&pageCount); err != nil {
		return 0, 0, fmt.Errorf("read fixture page count: %w", err)
	}
	if err := c.db.QueryRowContext(ctx, "PRAGMA max_page_count").Scan(&maxPageCount); err != nil {
		return 0, 0, fmt.Errorf("read fixture page limit: %w", err)
	}
	return pageCount, maxPageCount, nil
}

func (c *fixtureCapacityController) setMaxPageCount(ctx context.Context, value int64) error {
	if value <= 0 {
		return errors.New("fixture page limit must be positive")
	}
	if _, err := c.db.ExecContext(ctx, fmt.Sprintf("PRAGMA max_page_count = %d", value)); err != nil {
		return fmt.Errorf("set fixture page limit: %w", err)
	}
	return nil
}
