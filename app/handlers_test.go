package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestServer() http.Handler {
	return (&Server{store: NewMemoryStore()}).routes()
}

func TestHealth(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest("GET", "/health", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestShortenThenRedirect(t *testing.T) {
	h := newTestServer()

	// 1. Shorten a link
	body := strings.NewReader(`{"url":"https://example.com/some/long/path"}`)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("POST", "/shorten", body))

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp shortenResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	if len(resp.Code) != codeLength {
		t.Fatalf("expected a %d-character code, got %q", codeLength, resp.Code)
	}

	// 2. Follow the short code
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/"+resp.Code, nil))

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}
	if got := rec.Header().Get("Location"); got != "https://example.com/some/long/path" {
		t.Fatalf("redirected to the wrong place: %q", got)
	}
}

func TestShortenRejectsBadURL(t *testing.T) {
	for _, bad := range []string{`{"url":"not a url"}`, `{"url":"ftp://x.com"}`, `{}`, `oops`} {
		rec := httptest.NewRecorder()
		newTestServer().ServeHTTP(rec, httptest.NewRequest("POST", "/shorten", strings.NewReader(bad)))
		if rec.Code != http.StatusBadRequest {
			t.Errorf("input %s: expected 400, got %d", bad, rec.Code)
		}
	}
}

func TestUnknownCodeReturns404(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestServer().ServeHTTP(rec, httptest.NewRequest("GET", "/nope12", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
