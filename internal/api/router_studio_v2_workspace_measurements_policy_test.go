package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

const measurementsPath = "/api/v1/datalink/studio-v2/workspace/measurements"

func measurementRequest(t *testing.T, router *gin.Engine, method, path, body string) (status int, decoded map[string]any) {
	t.Helper()
	req := httptest.NewRequestWithContext(context.Background(), method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &decoded), recorder.Body.String())
	return recorder.Code, decoded
}

func counterBody(policy string) string {
	return `{"id":"meas-energy","device_id":"dev-1","point_id":"pt-1","name":"Energy","quantity":"energy","unit":"kWh",
		"semantic_kind":"counter","numeric_encoding":"uint16","counter_policy":` + policy + `}`
}

// TestStudioV2WorkspaceMeasurements_ExactRolloverPolicy drives the exact rollover
// evidence through the real route: it is saved and returned as the text it was
// sent as, a JSON number is refused instead of coerced, an unusable value is
// refused with a static message, and a changed policy moves the series epoch.
func TestStudioV2WorkspaceMeasurements_ExactRolloverPolicy(t *testing.T) {
	router := newStudioV2WorkspaceMeasurementRouter(t)

	code, created := measurementRequest(t, router, http.MethodPost, measurementsPath,
		counterBody(`{"allow_negative":false,"rollover_modulus":"9007199254740993","max_rate_per_second":"0.125"}`))
	require.Equal(t, http.StatusCreated, code, created)
	policy := created["data"].(map[string]any)["counter_policy"].(map[string]any)
	require.Equal(t, "9007199254740993", policy["rollover_modulus"], "exact above 2^53 as text")
	require.Equal(t, "0.125", policy["max_rate_per_second"])

	rejected := map[string]string{
		"a number for the modulus": `{"allow_negative":false,"rollover_modulus":65536}`,
		"a number for the rate":    `{"allow_negative":false,"rollover_modulus":"65536","max_rate_per_second":100}`,
	}
	for name, body := range rejected {
		fresh := newStudioV2WorkspaceMeasurementRouter(t)
		code, response := measurementRequest(t, fresh, http.MethodPost, measurementsPath, counterBody(body))
		require.GreaterOrEqual(t, code, 400, name)
		require.Less(t, code, 500, name)
		require.Equal(t, false, response["success"], name)
		_, listed := measurementRequest(t, fresh, http.MethodGet, measurementsPath, "")
		require.NotContains(t, fmt.Sprint(listed), "meas-energy", "nothing was saved after a refused "+name)
	}
	for name, body := range map[string]string{
		"a zero modulus":    `{"allow_negative":false,"rollover_modulus":"0"}`,
		"an exponent":       `{"allow_negative":false,"rollover_modulus":"6.5536e4"}`,
		"a rate alone":      `{"allow_negative":false,"max_rate_per_second":"100"}`,
		"negatives allowed": `{"allow_negative":true,"rollover_modulus":"65536"}`,
	} {
		fresh := newStudioV2WorkspaceMeasurementRouter(t)
		code, response := measurementRequest(t, fresh, http.MethodPost, measurementsPath, counterBody(body))
		require.Equal(t, http.StatusUnprocessableEntity, code, name)
		failure := response["error"].(map[string]any)
		require.Equal(t, "MEASUREMENT_CREATE_FAILED", failure["code"], name)
		require.NotContains(t, failure["message"], "6.5536e4", "the refusal is a static message, not an echo of the value")
	}

	update := func(policy string) (bool, int) {
		body := strings.Replace(counterBody(policy), `"id":"meas-energy",`, "", 1)
		code, response := measurementRequest(t, router, http.MethodPut, measurementsPath+"/meas-energy", body)
		if code != http.StatusOK {
			return false, code
		}
		return response["data"].(map[string]any)["epoch_transitioned"].(bool), code
	}
	moved, code := update(`{"allow_negative":false,"rollover_modulus":"9007199254740993.00","max_rate_per_second":"0.1250"}`)
	require.Equal(t, http.StatusOK, code)
	require.False(t, moved, "the same numbers with another scale keep the series")
	moved, code = update(`{"allow_negative":false,"rollover_modulus":"65536","max_rate_per_second":"100"}`)
	require.Equal(t, http.StatusOK, code)
	require.True(t, moved, "changed evidence is a new series epoch")
	_, code = update(`{"allow_negative":false,"rollover_modulus":"-1"}`)
	require.Equal(t, http.StatusUnprocessableEntity, code, "an invalid revision is refused and the saved policy stays")
}
