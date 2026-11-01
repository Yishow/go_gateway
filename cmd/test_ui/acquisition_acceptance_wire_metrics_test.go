package main

import (
	"context"
	"sync"
	"time"
)

// Test-local application-frame accounting. It does not count TCP packets,
// Ethernet bytes or physical serial timing. Delay is real wall-clock waiting.
type acquisitionAcceptanceWireMetrics struct {
	mu                                                      sync.Mutex
	rtt                                                     time.Duration
	dials, requests, responses, requestBytes, responseBytes uint64
}

type acquisitionAcceptanceWireSnapshot struct {
	Dials, Requests, Responses, RequestBytes, ResponseBytes uint64
}

func (p *acquisitionWirePeer) acceptanceMetrics(rtt time.Duration) *acquisitionAcceptanceWireMetrics {
	p.mu.Lock()
	defer p.mu.Unlock()
	m := &acquisitionAcceptanceWireMetrics{rtt: rtt}
	p.measurement = m
	return m
}

func (p *acquisitionWirePeer) acceptanceMeasurement() *acquisitionAcceptanceWireMetrics {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.measurement
}

func (p *acquisitionWirePeer) acceptanceConnected() {
	if m := p.acceptanceMeasurement(); m != nil {
		m.mu.Lock()
		m.dials++
		m.mu.Unlock()
	}
}

func (p *acquisitionWirePeer) acceptanceDelay(ctx context.Context) bool {
	m := p.acceptanceMeasurement()
	if m == nil || m.rtt <= 0 {
		return true
	}
	timer := time.NewTimer(m.rtt)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func (p *acquisitionWirePeer) acceptanceRequest(request int) {
	if m := p.acceptanceMeasurement(); m != nil {
		m.mu.Lock()
		m.requests++
		m.requestBytes += uint64(request)
		m.mu.Unlock()
	}
}

func (p *acquisitionWirePeer) acceptanceResponse(response int) {
	if m := p.acceptanceMeasurement(); m != nil {
		m.mu.Lock()
		m.responses++
		m.responseBytes += uint64(response)
		m.mu.Unlock()
	}
}

func (m *acquisitionAcceptanceWireMetrics) snapshot() acquisitionAcceptanceWireSnapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	return acquisitionAcceptanceWireSnapshot{m.dials, m.requests, m.responses, m.requestBytes, m.responseBytes}
}
