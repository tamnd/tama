// Package server is the HTTP layer: the JSON API under /api and the embedded
// web app for everything else, with SPA fallback to index.html.
//
// Router: Go 1.22+ stdlib http.ServeMux with method and wildcard patterns.
// chi comes in only if a concrete middleware need appears; the swap is one
// line in New.
package server

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/tamnd/tama/pkg/api"
	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/course"
	"github.com/tamnd/tama/pkg/store"
	"github.com/tamnd/tama/web"
)

// devUpstream is where `tama serve --dev` proxies non-API requests: the Vite
// dev server.
const devUpstream = "http://127.0.0.1:5173"

// Options are the serve-time knobs New cannot read from config.
type Options struct {
	// Version renders in /api/healthz.
	Version string
	// Dev swaps the embedded dist for a reverse proxy to the Vite server.
	Dev bool
}

// Server holds the pieces the handlers need.
type Server struct {
	cfg     *config.Config
	db      *store.DB
	opts    Options
	handler http.Handler

	ready chan struct{}
	bound string
}

// New wires the routes and the middleware chain, outermost first: request
// ID, real IP, request logging, panic recovery, session auth.
func New(cfg *config.Config, db *store.DB, opts Options) *Server {
	s := &Server{cfg: cfg, db: db, opts: opts, ready: make(chan struct{})}

	h := api.New(db)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", s.handleHealthz)
	h.Routes(mux)
	mux.HandleFunc("GET /api/languages", h.RequireUser(s.handleLanguages))
	mux.HandleFunc("GET /api/catalog", h.RequireUser(s.handleCatalog))
	if opts.Dev {
		mux.Handle("/", devProxy())
	} else {
		mux.Handle("/", spaHandler())
	}

	s.handler = requestID(realIP(logRequests(recovered(h.SessionMiddleware(mux)))))
	return s
}

// Run binds, announces readiness, and serves until the context is cancelled,
// then drains in-flight requests with a 10 second deadline. The WAL
// checkpoint happens in store.Close, which the caller owns.
func (s *Server) Run(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err // net's error already names the addr
	}
	s.bound = ln.Addr().String()
	close(s.ready)

	srv := &http.Server{
		Handler:           s.handler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	sweepCtx, stopSweeper := context.WithCancel(ctx)
	defer stopSweeper()
	go s.sweepSessions(sweepCtx)

	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

// Ready closes once the listener is bound; BoundAddr is valid after that.
// Tests binding port 0 use the pair to find the real port.
func (s *Server) Ready() <-chan struct{} {
	return s.ready
}

// BoundAddr is the address the listener actually got.
func (s *Server) BoundAddr() string {
	return s.bound
}

// sweepSessions deletes expired sessions once an hour until ctx ends.
func (s *Server) sweepSessions(ctx context.Context) {
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			n, err := s.db.DeleteExpiredSessions(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "session sweep failed", "err", err)
				continue
			}
			slog.DebugContext(ctx, "session sweep", "deleted", n)
		}
	}
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	api.WriteData(w, map[string]string{"status": "ok", "version": s.opts.Version})
}

func (s *Server) handleLanguages(w http.ResponseWriter, r *http.Request) {
	api.WriteData(w, course.Languages)
}

// handleCatalog lists seed courses for a base language, ?from=en by default.
func (s *Server) handleCatalog(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("from")
	if base == "" {
		base = "en"
	}
	courses := course.SeedCatalog(base)
	if courses == nil {
		api.WriteError(w, api.CodeNotFound, "no such base language")
		return
	}
	api.WriteData(w, courses)
}

// spaHandler serves the embedded web build. Hashed assets get immutable
// cache headers; index.html is no-cache so a new deploy shows up on refresh;
// unknown paths fall back to index.html so client-side routes deep-link.
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
				if strings.HasPrefix(p, "assets/") {
					// Vite content-hashes everything under assets/.
					w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
				}
				fileServer.ServeHTTP(w, r)
				return
			} else if !errors.Is(err, fs.ErrNotExist) {
				api.WriteError(w, api.CodeInternal, "internal error")
				return
			}
		}
		w.Header().Set("Cache-Control", "no-cache")
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}

// devProxy forwards non-API traffic to the Vite dev server so Go and the
// hot-reloading UI share one origin.
func devProxy() http.Handler {
	target, err := url.Parse(devUpstream)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(target)
}
