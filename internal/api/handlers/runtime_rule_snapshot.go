package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/datalink/workspace"
	"reflect"
)

var errRuntimeStaleReference = errors.New("stale-reference")

// A legacy saved-only read must not stitch two observed revisions together.
// This does not establish effective facts; task 1.4 supplies that authority.
func (h *RuntimeHandler) checkRuntimeSavedSnapshot(ctx context.Context, before *workspace.Record, rules []byte) error {
	after, err := h.workspaceSvc.GetOrCreate(ctx)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(before, after) {
		return errRuntimeStaleReference
	}
	latest, err := h.listWorkspaceSetupSourceRules(ctx, before.OrderedDeviceIDs)
	if err != nil {
		return err
	}
	if !bytes.Equal(rules, runtimeRuleSnapshot(latest)) {
		return errRuntimeStaleReference
	}
	return nil
}

func runtimeRuleSnapshot(rules []*schema.SourceRule) []byte {
	data, _ := json.Marshal(rules)
	return data
}

func (h *RuntimeHandler) checkRuntimeConnectorSnapshot(ctx context.Context, before *schema.DatabaseConnector) error {
	after, err := h.dbConnectorReader.GetByID(ctx, before.ID)
	if err != nil {
		return err
	}
	if after == nil || before.ID != after.ID || before.IdentityRevision != after.IdentityRevision || before.ConnectionConfig != after.ConnectionConfig || before.Name != after.Name || before.Kind != after.Kind || before.DefaultWriteIntervalSeconds != after.DefaultWriteIntervalSeconds {
		return errRuntimeStaleReference
	}
	return nil
}
