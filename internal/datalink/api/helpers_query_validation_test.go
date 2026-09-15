package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryIntegerFallsBackOnInvalidOrOverflow(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  int
	}{
		{"", 25}, {"bad", 25}, {"-1", 25}, {"1.5", 25},
		{"999999999999999999999999999", 25}, {"0", 0}, {"12", 12},
	} {
		t.Run(tc.value, func(t *testing.T) {
			request := httptest.NewRequestWithContext(t.Context(), "GET", "/?limit="+tc.value, http.NoBody)
			require.Equal(t, tc.want, getQueryParamInt(request, "limit", 25))
		})
	}
}
