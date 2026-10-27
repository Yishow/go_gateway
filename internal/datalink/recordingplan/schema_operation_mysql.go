package recordingplan

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
)

const (
	mysqlSchemaApplyManifestVersion = 1
	mysqlSchemaStepPending          = "pending"
	mysqlSchemaStepStarted          = "started"
	mysqlSchemaStepVerified         = "verified"
	mysqlSchemaStepFailed           = "failed"
)

var (
	errMySQLSchemaApplyManifest = errors.New("mysql schema apply manifest is invalid")
	mysqlBaselineVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:[-+~:].*)?$`)
)

type mysqlSchemaApplyManifest struct {
	Version     int                    `json:"version"`
	TokenDigest string                 `json:"token_digest"`
	ScopeKey    string                 `json:"scope_key"`
	Steps       []mysqlSchemaApplyStep `json:"steps"`
}

type mysqlSchemaApplyStep struct {
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	SQLHash     string `json:"sql_hash"`
	Status      string `json:"status"`
	Fingerprint string `json:"fingerprint,omitempty"`
	FailureCode string `json:"failure_code,omitempty"`
}

type mysqlSchemaOperationResumeRepository interface {
	ReclaimMySQLSchemaApplyOperation(ctx context.Context, operationID, workspaceID, token, scopeKey, expectedDetail, newOwner string) (*SchemaOperation, ClaimOutcome, error)
}

func newMySQLSchemaApplyManifest(token *SchemaPreviewToken) (mysqlSchemaApplyManifest, error) {
	steps, err := mysqlSchemaApplySteps(token)
	if err != nil {
		return mysqlSchemaApplyManifest{}, err
	}
	scopeKey := schemaOperationScopeKey(token)
	if scopeKey == "" || token.Digest == "" {
		return mysqlSchemaApplyManifest{}, errMySQLSchemaApplyManifest
	}
	return mysqlSchemaApplyManifest{
		Version: mysqlSchemaApplyManifestVersion, TokenDigest: token.Digest, ScopeKey: scopeKey, Steps: steps,
	}, nil
}

func mysqlSchemaApplySteps(token *SchemaPreviewToken) ([]mysqlSchemaApplyStep, error) {
	if token == nil || token.GroupLayout == nil || token.GroupLayout.MySQL == nil || token.Dialect != dialectMySQL {
		return nil, ErrSchemaAdapterUnverified
	}
	canonical, err := GenerateGroupSchemaDDL(dialectMySQL, *token.GroupLayout)
	if err != nil {
		return nil, err
	}
	steps := make([]mysqlSchemaApplyStep, 0, len(token.Statements))
	previous := -1
	for _, statement := range token.Statements {
		index := slices.Index(canonical, statement)
		if index < 0 || index <= previous {
			return nil, fmt.Errorf("%w: statement sequence", ErrPreviewTokenStale)
		}
		kind, name, err := mysqlSchemaStepIdentity(*token.GroupLayout, index)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256([]byte(statement))
		steps = append(steps, mysqlSchemaApplyStep{
			Kind: kind, Name: name, SQLHash: hex.EncodeToString(sum[:]), Status: mysqlSchemaStepPending,
		})
		previous = index
	}
	return steps, nil
}

func mysqlSchemaStepIdentity(layout GroupSchemaLayout, statementIndex int) (kind, name string, err error) {
	switch {
	case statementIndex == 0:
		return mysqlStepKindGroupTable, layout.TableName, nil
	case statementIndex > 0 && statementIndex <= len(layout.Indexes):
		return mysqlStepKindIndex, layout.Indexes[statementIndex-1].Name, nil
	case statementIndex == len(layout.Indexes)+1:
		return mysqlStepKindReceiptTable, managedEffectReceiptTable, nil
	default:
		return "", "", fmt.Errorf("%w: MySQL DDL step", ErrPreviewTokenStale)
	}
}

func mysqlSchemaApplyManifestFromDetail(detail string, token *SchemaPreviewToken) (mysqlSchemaApplyManifest, error) {
	var manifest mysqlSchemaApplyManifest
	if err := json.Unmarshal([]byte(detail), &manifest); err != nil {
		return mysqlSchemaApplyManifest{}, fmt.Errorf("%w: decode detail", errMySQLSchemaApplyManifest)
	}
	expected, err := newMySQLSchemaApplyManifest(token)
	if err != nil {
		return mysqlSchemaApplyManifest{}, err
	}
	if manifest.Version != expected.Version || manifest.ScopeKey != expected.ScopeKey ||
		subtle.ConstantTimeCompare([]byte(manifest.TokenDigest), []byte(expected.TokenDigest)) != 1 ||
		len(manifest.Steps) != len(expected.Steps) {
		return mysqlSchemaApplyManifest{}, errMySQLSchemaApplyManifest
	}
	for index, step := range manifest.Steps {
		want := expected.Steps[index]
		if step.Kind != want.Kind || step.Name != want.Name || step.SQLHash != want.SQLHash || !validMySQLSchemaStepStatus(step.Status) {
			return mysqlSchemaApplyManifest{}, errMySQLSchemaApplyManifest
		}
		if step.Status == mysqlSchemaStepVerified && step.Fingerprint == "" {
			return mysqlSchemaApplyManifest{}, errMySQLSchemaApplyManifest
		}
		if step.Status != mysqlSchemaStepVerified && step.Fingerprint != "" {
			return mysqlSchemaApplyManifest{}, errMySQLSchemaApplyManifest
		}
	}
	return manifest, nil
}

func validMySQLSchemaStepStatus(status string) bool {
	switch status {
	case mysqlSchemaStepPending, mysqlSchemaStepStarted, mysqlSchemaStepVerified, mysqlSchemaStepFailed:
		return true
	default:
		return false
	}
}

func encodeMySQLSchemaApplyManifest(manifest mysqlSchemaApplyManifest) (string, error) {
	encoded, err := json.Marshal(manifest)
	if err != nil {
		return "", fmt.Errorf("encode MySQL schema apply manifest: %w", err)
	}
	return string(encoded), nil
}

func mysqlSchemaApplyBaselineSupported(layout *GroupSchemaLayout) bool {
	if layout == nil || layout.MySQL == nil || layout.Partition != nil {
		return false
	}
	options := layout.MySQL
	matches := mysqlBaselineVersionPattern.FindStringSubmatch(options.ServerVersion)
	if len(matches) != 4 {
		return false
	}
	switch options.ServerFamily {
	case mysqlFamilyMySQL:
		return matches[1] == "8" && matches[2] == "4"
	case mysqlFamilyMariaDB:
		return matches[1] == "11" && matches[2] == "4"
	default:
		return false
	}
}
