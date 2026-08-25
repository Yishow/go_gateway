package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
)

func newHandlerTestRequest(method, target string, body io.Reader) *http.Request {
	return httptest.NewRequestWithContext(context.Background(), method, target, body)
}
