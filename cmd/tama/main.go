// Command tama runs the tama language learning server and its content
// tooling. The app is a single binary: the web UI is embedded, the database
// is one SQLite file, and course packs live next to it in the data directory.
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tamnd/tama/pkg/cli"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	os.Exit(cli.Execute(ctx))
}
