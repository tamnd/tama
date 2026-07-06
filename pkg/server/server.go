// Package server is the HTTP layer: the JSON API under /api and the embedded
// web app for everything else, with SPA fallback to index.html.
package server

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/course"
	"github.com/tamnd/tama/pkg/store"
	"github.com/tamnd/tama/web"
)

// Server holds the pieces the handlers need.
type Server struct {
	cfg *config.Config
	st  *store.DB
	mux *http.ServeMux
}

// New wires the routes.
func New(cfg *config.Config, st *store.DB) *Server {
	s := &Server{cfg: cfg, st: st, mux: http.NewServeMux()}

	s.mux.HandleFunc("GET /api/health", s.handleHealth)
	s.mux.HandleFunc("GET /api/languages", s.handleLanguages)
	s.mux.HandleFunc("GET /api/catalog", s.handleCatalog)
	s.mux.Handle("/", spaHandler())

	return s
}

// Run serves until the context is cancelled, then shuts down gracefully.
func (s *Server) Run(ctx context.Context) error {
	srv := &http.Server{
		Addr:              s.cfg.Addr,
		Handler:           s.mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleLanguages(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, course.Languages)
}

// handleCatalog lists seed courses for a base language, ?from=en by default.
func (s *Server) handleCatalog(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("from")
	if base == "" {
		base = "en"
	}
	courses := course.SeedCatalog(base)
	if courses == nil {
		writeError(w, http.StatusNotFound, "unknown_base", "no such base language")
		return
	}
	writeJSON(w, http.StatusOK, courses)
}

// spaHandler serves the embedded web build. Unknown paths fall back to
// index.html so client-side routes deep-link.
func spaHandler() http.Handler {
	dist, err := fs.Sub(web.Dist, "dist")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(dist))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/")
		if p != "" {
			if _, err := fs.Stat(dist, p); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			} else if !errors.Is(err, fs.ErrNotExist) {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, errCode, msg string) {
	writeJSON(w, code, map[string]any{"error": map[string]string{"code": errCode, "message": msg}})
}
