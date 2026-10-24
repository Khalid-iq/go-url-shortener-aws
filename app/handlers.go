package main

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/url"
)

//go:embed static/index.html
var staticFiles embed.FS

const codeChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const codeLength = 6

type Server struct {
	store Store
}

// routes connects URL patterns to handler functions.
// Go 1.22+ supports the method ("GET") and path values ({code}) in patterns.
func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.handleHome) // {$} means exactly "/", nothing after it
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /shorten", s.handleShorten)
	mux.HandleFunc("GET /{code}", s.handleRedirect)
	return mux
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	page, _ := staticFiles.ReadFile("static/index.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}

// handleHealth is what the load balancer will call in Stage 6
// to check the app is alive.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func (s *Server) handleShorten(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Send JSON like {\"url\": \"https://example.com\"}")
		return
	}
	if !isValidURL(req.URL) {
		writeError(w, http.StatusBadRequest, "Enter a full link starting with http:// or https://")
		return
	}

	code, err := s.newCode()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Could not create a short code")
		return
	}
	if err := s.store.Save(code, req.URL); err != nil {
		writeError(w, http.StatusInternalServerError, "Could not save the link")
		return
	}

	writeJSON(w, http.StatusCreated, shortenResponse{
		Code:     code,
		ShortURL: baseURL(r) + "/" + code,
	})
}

func (s *Server) handleRedirect(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	target, err := s.store.Get(code)
	if errors.Is(err, ErrNotFound) {
		http.Error(w, "This short link doesn't exist.", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	// 302 = temporary redirect. Browsers don't cache it, so every visit hits
	// your app (useful later if you add click counting).
	http.Redirect(w, r, target, http.StatusFound)
}

// newCode generates a random code, retrying if it's already taken.
func (s *Server) newCode() (string, error) {
	for i := 0; i < 5; i++ {
		code, err := randomCode()
		if err != nil {
			return "", err
		}
		if !s.store.Exists(code) {
			return code, nil
		}
	}
	return "", errors.New("too many collisions")
}

func randomCode() (string, error) {
	b := make([]byte, codeLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeChars))))
		if err != nil {
			return "", err
		}
		b[i] = codeChars[n.Int64()]
	}
	return string(b), nil
}

func isValidURL(raw string) bool {
	u, err := url.ParseRequestURI(raw)
	if err != nil {
		return false
	}
	return (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

// baseURL builds "https://yourdomain.com" from the request.
// Behind a load balancer, the original scheme arrives in X-Forwarded-Proto.
func baseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
