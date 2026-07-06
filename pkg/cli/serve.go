package cli

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/spf13/cobra"

	"github.com/tamnd/tama/pkg/config"
	"github.com/tamnd/tama/pkg/server"
	"github.com/tamnd/tama/pkg/store"
)

func newServeCmd() *cobra.Command {
	var (
		dev         bool
		printConfig bool
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the tama server",
		Long: "Starts the HTTP server with the embedded web UI. State goes to a single\n" +
			"SQLite file under the data directory, course packs inside it.\n\n" +
			"With --dev, non-API requests proxy to the Vite dev server on\n" +
			"http://127.0.0.1:5173 instead of the embedded bundle, and the listen\n" +
			"address must stay on loopback.",
		Example: "  tama serve\n" +
			"  tama serve --addr 127.0.0.1:8080 --data /srv/tama\n" +
			"  tama serve --print-config",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load(cmd.Flags())
			if err != nil {
				return err
			}
			if printConfig {
				cfg.Print(cmd.OutOrStdout())
				return nil
			}
			setupLogging(cfg)
			for _, w := range cfg.Warnings {
				slog.Warn(w)
			}

			if dev {
				if cfg.Addr, err = loopbackAddr(cfg.Addr); err != nil {
					return err
				}
			}

			st, err := store.Open(cmd.Context(), cfg.DBPath())
			if err != nil {
				return fmt.Errorf("open store: %w", err)
			}
			defer st.Close()

			webSource := "embedded"
			if dev {
				webSource = "proxy"
			}
			srv := server.New(cfg, st, server.Options{Version: Version, Dev: dev})
			slog.Info("tama serving",
				"addr", cfg.Addr, "data", cfg.DataDir, "db", cfg.DBPath(), "web", webSource)
			return srv.Run(cmd.Context())
		},
	}

	cmd.Flags().String("addr", ":4321", "listen address")
	cmd.Flags().String("data", "", "data directory (default ~/.tama)")
	cmd.Flags().BoolVar(&dev, "dev", false, "proxy the UI to the Vite dev server (loopback only)")
	cmd.Flags().BoolVar(&printConfig, "print-config", false, "print the resolved config with sources and exit")
	return cmd
}

// setupLogging installs the default slog handler per config: text unless the
// config asks for JSON, wrapped so request IDs ride along on every line.
func setupLogging(cfg *config.Config) {
	level := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}[cfg.Log.Level]
	opts := &slog.HandlerOptions{Level: level}

	var h slog.Handler
	if cfg.Log.Format == "json" {
		h = slog.NewJSONHandler(os.Stderr, opts)
	} else {
		h = slog.NewTextHandler(os.Stderr, opts)
	}
	slog.SetDefault(slog.New(server.ContextHandler{Handler: h}))
}

// loopbackAddr enforces that dev mode never faces a network. An unspecified
// host narrows to 127.0.0.1; anything non-loopback is refused.
func loopbackAddr(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	if host == "" {
		return net.JoinHostPort("127.0.0.1", port), nil
	}
	if host == "localhost" {
		return addr, nil
	}
	if ip := net.ParseIP(host); ip != nil && ip.IsLoopback() {
		return addr, nil
	}
	return "", fmt.Errorf("--dev only binds loopback addresses, not %q", addr)
}
