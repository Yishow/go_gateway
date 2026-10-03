package recordingplan

import (
	"encoding/json"
	"fmt"
)

// Group fields share the existing tables JSON column. Legacy rows remain a
// plain array, while group rows use this additive envelope so no token table
// or migration is needed.
type persistedPreviewTables struct {
	Tables         []SchemaPreviewTable `json:"tables"`
	GroupLayout    *GroupSchemaLayout   `json:"group_layout,omitempty"`
	SourceDigest   string               `json:"source_digest,omitempty"`
	SchemaRevision string               `json:"schema_revision,omitempty"`
	SchemaDigest   string               `json:"schema_digest,omitempty"`
}

func encodePreviewTables(token *SchemaPreviewToken) ([]byte, error) {
	if token.GroupLayout == nil && token.SourceDigest == "" && token.SchemaRevision == "" && token.SchemaDigest == "" {
		return json.Marshal(token.Tables)
	}
	return json.Marshal(persistedPreviewTables{
		Tables: token.Tables, GroupLayout: token.GroupLayout, SourceDigest: token.SourceDigest,
		SchemaRevision: token.SchemaRevision, SchemaDigest: token.SchemaDigest,
	})
}

func decodePreviewTables(raw string, token *SchemaPreviewToken) error {
	var tables []SchemaPreviewTable
	if err := json.Unmarshal([]byte(raw), &tables); err == nil {
		token.Tables = tables
		return nil
	}
	var envelope persistedPreviewTables
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return fmt.Errorf("decode preview tables: %w", err)
	}
	if envelope.GroupLayout != nil {
		layout, err := ValidateGroupSchemaLayout(*envelope.GroupLayout)
		if err != nil {
			return fmt.Errorf("decode group schema layout: %w", err)
		}
		envelope.GroupLayout = &layout
	}
	token.Tables = envelope.Tables
	token.GroupLayout = envelope.GroupLayout
	token.SourceDigest = envelope.SourceDigest
	token.SchemaRevision = envelope.SchemaRevision
	token.SchemaDigest = envelope.SchemaDigest
	return nil
}
