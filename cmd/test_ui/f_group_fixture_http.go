//go:build f_write_group_fixture

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"go-gateway/internal/datalink/grouppipeline"
	"go-gateway/internal/datalink/measurement"
	datalinkruntime "go-gateway/internal/datalink/runtime"
)

type fixtureHTTPServer struct {
	server   *http.Server
	listener net.Listener
}

type fixtureClockRequest struct {
	At string `json:"at"`
}

type fixturePollRequest struct {
	GroupID       string   `json:"group_id"`
	PointIDs      []string `json:"point_ids,omitempty"`
	AcquisitionID string   `json:"acquisition_id,omitempty"`
}

type fixtureReleaseRequest struct {
	Indices []int `json:"indices"`
}

type fixtureCaptureResponse struct {
	Index      int                        `json:"index"`
	Status     string                     `json:"status"`
	Attempts   int                        `json:"attempts"`
	LastReason string                     `json:"last_reason,omitempty"`
	Envelope   measurement.SampleEnvelope `json:"envelope"`
}

type fixtureStateResponse struct {
	Now             string                      `json:"now"`
	Paused          bool                        `json:"paused"`
	CaptureCapacity int                         `json:"capture_capacity"`
	CaptureCount    int                         `json:"capture_count"`
	Captured        []fixtureCaptureResponse    `json:"captured"`
	Pipeline        []grouppipeline.GroupStatus `json:"pipeline"`
	PipelineErrors  uint64                      `json:"pipeline_errors"`
	Faults          []fixtureFaultState         `json:"faults"`
}

type fixturePollResponse struct {
	Indices     []int `json:"indices"`
	ResultCount int   `json:"result_count"`
}

type fixtureReleaseResult struct {
	Index    int    `json:"index"`
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason,omitempty"`
	Outcome  string `json:"outcome"`
}

type fixtureReleaseResponse struct {
	Results []fixtureReleaseResult `json:"results"`
}

type fixtureFaultsResponse struct {
	Faults []fixtureFaultState `json:"faults"`
}

func (f *groupFixture) startHTTP(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /state", f.handleState)
	mux.HandleFunc("GET /faults", f.handleFaults)
	mux.HandleFunc("POST /clock", f.handleClock)
	mux.HandleFunc("POST /pause", f.handlePause)
	mux.HandleFunc("POST /poll", f.handlePoll)
	mux.HandleFunc("POST /release", f.handleRelease)
	mux.HandleFunc("POST /tick", f.handleTick)
	mux.HandleFunc("POST /fault", f.handleFault)
	mux.HandleFunc("POST /capacity", f.handleCapacity)
	listenerConfig := net.ListenConfig{}
	listener, err := listenerConfig.Listen(ctx, "tcp", fmt.Sprintf("127.0.0.1:%d", f.port))
	if err != nil {
		return errors.New("fixture controller could not bind loopback")
	}
	f.httpMu.Lock()
	f.http = &fixtureHTTPServer{
		listener: listener,
		server:   &http.Server{Handler: mux, ReadHeaderTimeout: 2 * time.Second},
	}
	server := f.http.server
	f.httpMu.Unlock()
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			f.recordError(err)
		}
	}()
	return nil
}

func (f *groupFixture) stopHTTP(ctx context.Context) error {
	f.httpMu.Lock()
	server := f.http
	f.http = nil
	f.httpMu.Unlock()
	if server == nil {
		return nil
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := server.server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("stop fixture HTTP server: %w", err)
	}
	return nil
}

func (f *groupFixture) handleState(w http.ResponseWriter, _ *http.Request) {
	f.writeJSON(w, http.StatusOK, f.state())
}

func (f *groupFixture) handleFaults(w http.ResponseWriter, _ *http.Request) {
	var faults []fixtureFaultState
	if f.faults != nil {
		faults = f.faults.snapshot()
	}
	f.writeJSON(w, http.StatusOK, fixtureFaultsResponse{Faults: faults})
}

func (f *groupFixture) handleClock(w http.ResponseWriter, r *http.Request) {
	var request fixtureClockRequest
	if err := decodeFixtureJSON(r, &request); err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid clock request")
		return
	}
	value, err := parseFixtureTime(request.At, "")
	if err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid UTC clock")
		return
	}
	f.opMu.Lock()
	f.clock.Set(value)
	f.opMu.Unlock()
	f.writeJSON(w, http.StatusOK, f.state())
}

func (f *groupFixture) handlePause(w http.ResponseWriter, r *http.Request) {
	var request struct{}
	if err := decodeFixtureJSON(r, &request); err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid pause request")
		return
	}
	if err := f.pause(r.Context()); err != nil {
		f.writeError(w, http.StatusInternalServerError, "pause failed")
		return
	}
	f.writeJSON(w, http.StatusOK, f.state())
}

func (f *groupFixture) handlePoll(w http.ResponseWriter, r *http.Request) {
	var request fixturePollRequest
	if err := decodeFixtureJSON(r, &request); err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid poll request")
		return
	}
	response, err := f.poll(r.Context(), request)
	if err != nil {
		f.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	f.writeJSON(w, http.StatusOK, response)
}

func (f *groupFixture) handleRelease(w http.ResponseWriter, r *http.Request) {
	var request fixtureReleaseRequest
	if err := decodeFixtureJSON(r, &request); err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid release request")
		return
	}
	response, err := f.release(r.Context(), request.Indices)
	if err != nil {
		f.writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	f.writeJSON(w, http.StatusOK, response)
}

func (f *groupFixture) handleTick(w http.ResponseWriter, r *http.Request) {
	var request struct{}
	if err := decodeFixtureJSON(r, &request); err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid tick request")
		return
	}
	if err := f.tick(r.Context()); err != nil {
		f.writeError(w, http.StatusInternalServerError, "tick reconcile failed")
		return
	}
	f.writeJSON(w, http.StatusOK, f.state())
}

func (f *groupFixture) handleFault(w http.ResponseWriter, r *http.Request) {
	var request fixtureFaultRequest
	if err := decodeFixtureJSON(r, &request); err != nil {
		f.writeError(w, http.StatusBadRequest, "invalid fault request")
		return
	}
	state, err := f.faults.configure(r.Context(), request)
	if err != nil {
		f.writeError(w, http.StatusBadRequest, "fixture fault rejected")
		return
	}
	f.writeJSON(w, http.StatusOK, map[string]fixtureFaultState{"fault": state})
}

func (f *groupFixture) handleCapacity(w http.ResponseWriter, r *http.Request) {
	var request fixtureCapacityRequest
	if err := decodeFixtureJSON(r, &request); err != nil || request.Kind != fixtureCapacityKindDiskFull {
		f.writeError(w, http.StatusBadRequest, "invalid capacity request")
		return
	}
	if f.capacity == nil {
		f.writeError(w, http.StatusBadRequest, "fixture capacity controller unavailable")
		return
	}
	state, err := f.capacity.configure(r.Context(), request.Enabled)
	if err != nil {
		f.writeError(w, http.StatusBadRequest, "fixture capacity rejected")
		return
	}
	f.writeJSON(w, http.StatusOK, state)
}

func (f *groupFixture) state() fixtureStateResponse {
	f.capture.mu.Lock()
	captures := make([]fixtureCaptureResponse, 0, len(f.capture.samples))
	for index, sample := range f.capture.samples {
		captures = append(captures, fixtureCaptureResponse{
			Index: index, Status: sample.status, Attempts: sample.attempts,
			LastReason: sample.lastReason, Envelope: sample.envelope,
		})
	}
	capacity := f.capture.capacity
	f.capture.mu.Unlock()
	var pipeline []grouppipeline.GroupStatus
	if f.pipe != nil {
		pipeline = f.pipe.Status()
	}
	var faults []fixtureFaultState
	if f.faults != nil {
		faults = f.faults.snapshot()
	}
	return fixtureStateResponse{
		Now: f.clock.Now().Format(time.RFC3339Nano), Paused: f.paused.Load(),
		CaptureCapacity: capacity, CaptureCount: len(captures), Captured: captures,
		Pipeline: pipeline, PipelineErrors: f.errors.Load(), Faults: faults,
	}
}

func (f *groupFixture) writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		f.recordError(err)
	}
}

func (f *groupFixture) writeError(w http.ResponseWriter, status int, reason string) {
	f.writeJSON(w, status, map[string]string{"error": reason})
}

func decodeFixtureJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return errors.New("fixture request must contain one JSON object")
	}
	return nil
}

func releaseReasons(err error) []string {
	if err == nil {
		return nil
	}
	var joined interface{ Unwrap() []error }
	if errors.As(err, &joined) {
		var reasons []string
		for _, inner := range joined.Unwrap() {
			reasons = append(reasons, releaseReasons(inner)...)
		}
		return reasons
	}
	var sampleError *datalinkruntime.GroupSampleError
	if errors.As(err, &sampleError) && strings.TrimSpace(sampleError.Reason) != "" {
		return []string{sampleError.Reason}
	}
	return []string{"release-failed"}
}
