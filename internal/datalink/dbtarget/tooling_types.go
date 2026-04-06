package dbtarget

import "time"

type SchemaGenerateRequest struct {
	DryRun bool `json:"dry_run"`
}

type SchemaGenerateResult struct {
	ConnectorID string   `json:"connector_id"`
	DryRun      bool     `json:"dry_run"`
	Statements  []string `json:"statements"`
	Executed    int      `json:"executed"`
}

type MappingDryRunRequest struct {
	CandidateIDs []string `json:"candidate_ids,omitempty"`
}

type MappingDryRunCandidateResult struct {
	CandidateID string `json:"candidate_id"`
	Status      string `json:"status"`
	Code        string `json:"code,omitempty"`
	Reason      string `json:"reason,omitempty"`
	MappingID   string `json:"mapping_id,omitempty"`
	TagID       string `json:"tag_id,omitempty"`
}

type MappingDryRunResult struct {
	ConnectorID string                         `json:"connector_id"`
	Results     []MappingDryRunCandidateResult `json:"results"`
}

type WriteHistoryRecord struct {
	Timestamp    time.Time `json:"timestamp"`
	Status       string    `json:"status"`
	RowCount     int       `json:"row_count"`
	ErrorSummary string    `json:"error_summary,omitempty"`
}
