package handlers

import (
	"net/http"
	"testing"
)

func TestWorkspaceRuleSwitchesDefaultScaleToMixed(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"plain","device_id":"mc-dev","start_address":"D39","count":1,"data_type":"int16","naming_prefix":"P","enabled":true,"scale_multiplier":1,"scale_offset":0}`
	if code, body := f.do(t, http.MethodPost, "/workspace/source-rules", create); code != http.StatusCreated {
		t.Fatalf("create rule: %d %v", code, body)
	}
	update := `{"device_id":"mc-dev","naming_prefix":"P","enabled":true,"layout_mode":"mixed","mixed_items":[{"item_id":"d39","address":"D39","data_type":"int16","semantic_kind":"gauge"}],"skipped_addresses":[],"share_enabled":false}`
	if code, body := f.do(t, http.MethodPut, "/workspace/source-rules/plain", update); code != http.StatusOK {
		t.Fatalf("switch rule to mixed: %d %v", code, body)
	}
}

func TestWorkspaceRuleSwitchDoesNotDiscardCustomScale(t *testing.T) {
	f := newMixedRuleHandlerFixture(t)
	create := `{"id":"plain","device_id":"mc-dev","start_address":"D39","count":1,"data_type":"int16","naming_prefix":"P","enabled":true,"scale_multiplier":2,"scale_offset":0}`
	if code, body := f.do(t, http.MethodPost, "/workspace/source-rules", create); code != http.StatusCreated {
		t.Fatalf("create rule: %d %v", code, body)
	}
	update := `{"layout_mode":"mixed","mixed_items":[{"item_id":"d39","address":"D39","data_type":"int16","semantic_kind":"gauge"}],"skipped_addresses":[],"share_enabled":false}`
	if code, body := f.do(t, http.MethodPut, "/workspace/source-rules/plain", update); code != http.StatusBadRequest {
		t.Fatalf("custom scale must not be discarded: %d %v", code, body)
	}
}
