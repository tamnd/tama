# tama

tama (タマ) is a self-hosted language learning app that ships as a single Go binary.
It gives you the full Duolingo-style experience: a path of units and lessons, hearts, streaks, XP, gems, leagues, quests, and a mascot cat cheering you on.
Any base language to any target language is a valid course.
Course content is generated ahead of time by an LLM through an OpenAI-compatible endpoint, validated, and frozen into course packs, so lessons never wait on a model.

Status: early. The foundation is being laid, expect rough edges everywhere.

## Run it

```sh
go build ./cmd/tama
./tama serve
```

Then open http://localhost:4321.

## Configuration

Everything is a flag, an environment variable, a config file entry, or a default, in that order.

| Env | Default | What |
| --- | --- | --- |
| `TAMA_ADDR` | `:4321` | listen address |
| `TAMA_DATA` | `~/.tama` | data directory (database, course packs) |
| `TAMA_LLM_BASE_URL` | `http://127.0.0.1:8000/v1` | OpenAI-compatible endpoint for course generation |
| `TAMA_LLM_API_KEY` | empty | API key for that endpoint |
| `TAMA_LLM_MODEL` | empty | model name to request |
| `TAMA_LOG_LEVEL` | `info` | `debug`, `info`, `warn`, or `error`; append `,json` for JSON logs |

The config file is `config.toml` in the data directory, so `~/.tama/config.toml` by default.
Every key is optional; unknown keys are warnings, not errors.

```toml
[server]
addr = ":4321"
data = "~/.tama"

[llm]
base_url = "http://127.0.0.1:8000/v1"
api_key = ""
model = ""
request_timeout = "5m"
connect_timeout = "30s"

[log]
level = "info"   # debug, info, warn, error
format = "text"  # text or json
```

`tama serve --print-config` prints the resolved values and where each one came from.

The server never calls the LLM during lessons.
Generation happens through `tama gen` or the admin queue, and the results are cached as course packs under the data directory.

## Layout

- `cmd/tama`: the binary.
- `pkg/`: server, store, course model, exercises, engine, generation.
- `web/`: the React app, built and embedded into the binary.

## License

MIT.
