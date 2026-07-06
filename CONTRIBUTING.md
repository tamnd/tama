# Contributing

## Pull requests

Work on a branch and open a PR; nothing lands on main directly.
Keep PRs focused on one slice of work, and keep commits small with imperative, capitalized titles ("Add session sweeper", not "added sweeper stuff").
Bodies stay plain and short: what changed and why, no ceremony.
Merge the open PR before stacking the next one on top of it.

Before you push:

```sh
go build ./...
go test ./... -race
go vet ./...
gofmt -s -l .   # must print nothing
```

The web app has its own checks under `web/`; run `make check` if you touched it.

## Layout rules

Everything lives under `pkg/`; there are no `internal/` directories, and a guard test fails the build if one appears.
Import direction is also pinned by test: `pkg/store` imports no siblings, `pkg/api` may reach `store`, `course`, `exercise`, `engine`, and `gen`, and nothing imports `cmd/`.
Only `pkg/config` reads the environment.

## Dependencies

Direct Go dependencies are allowlisted in `pkg/server/guard_test.go`.
Adding one means a row in that test, a row in this table, and a reason a reviewer can agree with.
Indirect modules pulled in by these (the charm stack behind fang, sqlite's build deps) are fine.

| Module | Why |
| --- | --- |
| modernc.org/sqlite | Pure-Go SQLite driver, keeps the binary CGO-free |
| github.com/spf13/cobra | Command tree for the CLI |
| github.com/spf13/pflag | Flag sets that config precedence is built on |
| github.com/charmbracelet/fang | Nicer help and errors around cobra |
| golang.org/x/crypto | argon2id for password hashing |
| github.com/BurntSushi/toml | config.toml parsing with unknown-key detection |
| github.com/klauspost/compress | zstd for course pack blobs |
