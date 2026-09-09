# AGENTS.md

Guidance for AI agents working in the `paepcke.de/osm2addr` repository.

> ## FIXED REQUIREMENT — EVERY CHANGE, NO EXCEPTIONS
>
> When a task or change is considered done, all five steps below MUST be completed
> in this exact order. 
>
> 1. **Format source code** — run `gofmt -w .` (or `make check`) so the tree
>    stays gofmt-clean.
> 2. **Build** — `go build -o osm2addr ./cmd/osm2addr` must succeed.
> 3. **Test** — `go test -count=1 ./...` must be green
> 4. **Commit** — `git add . && git commit -m '<message>'`.
> 5. **Tag** — bump the patch segment only: the result is `v0.0.<N+1>`. Never move, delete, or reuse an existing tag.
>
> These steps are non-negotiable for every single task regardless of size.
> Do not ask, just perform
>
> When build unit tests when possible, ensure unit test are run in parallel and fast.
> Use for unit test only small sample size data.

## Overview 

`osm2addr` parses OpenStreetMap PBF extract files and emits normalized,
postal-address mapping tables (JSON) keyed by a deterministic place ID. It is
optimized for the German (`DE`) address dataset but the architecture is
country-aware via the `uniform` dispatch. Output is intended to feed address
validation databases (MongoDB, etc.).

The module path is `paepcke.de/osm2addr`. The Go module was initialized with
`go mod init paepcke.de/osm2addr`; do not change this path. `go.mod` tracks
`go 1.26.5`.

## Running the Tool

`cmd/osm2addr/main.go` is the only entrypoint. CLI signature:

```
osm2addr [<COUNTRY_CODE> [<PATH_TO_PBF>]]
```

- `COUNTRY_CODE` must be exactly 2 characters (uppercase). Default `DE`.
- `PATH_TO_PBF` defaults to `data/germany-latest.osm.pbf` (relative to the
  process CWD). `main.go` opens the file exactly once, after argument parsing
  and the startup report.
- Output is always written to `json/<COUNTRY_CODE>/` (relative to the
  process CWD), not next to the binary.

### Preload CSV (optional but expected for `DE`)

If a file exists at
`<dirname of PBF>/validated-preload/<COUNTRY_CODE>.csv`, it is parsed as a
comma-separated CSV with **6 fields per record**, a header row (skipped), and
fixed column indices:

- `City`        = column **2**
- `Postcode`    = column **3**
- Expected postcode length = **5** (German)

Records failing length checks — and CSV rows that fail to parse (e.g. wrong
field count, stray quote) — are counted into
`json/<CC>/error.preload.json` instead of aborting the run. The preload
feeds the same `targets` channel
as the PBF parser, so its entries participate in dedup and place-ID
generation. If the preload file is missing, the run proceeds without it
(`OSM:PreLoadFile # n/a`).

### Trusted preload.csv (optional)

At startup, before the validated preload, `trustedPreloadFeed`
(`preload_trusted.go`) looks for a file named `preload.csv` next to the PBF
file (process CWD as fallback). The first line must be a valid UTF-8 CSV
header whose separator is auto-detected (`,` assumed by default; `;`, tab,
`|` detected) and which contains the exact field names `POSTLEITZAHL`,
`ORT_NAME` and `STRASSE_NAME`. If confirmed, the rows are fed into `targets`
**first** as trusted postcode/city/street values: `uniform` normalization is
deliberately skipped and no matching is attempted. Malformed rows are counted
into `json/<CC>/error.preload.trusted.json`. Order in `core.go::Parse`:
collector → trusted preload (wait) → validated preload (wait) → PBF parser.

## Architecture & Data Flow

The pipeline is a fixed three-stage goroutine chain wired in
`core.go::Parse`:

```
preload.csv (trusted) ─┐
preloadFeed (csv)  ───┼──>  chan *tagSet  targets  ──>  collect ──> JSON files
pbfparser (pbf)   ────┘
```

1. **`trustedPreloadFeed`** (`preload_trusted.go`) — started first, in its
   own goroutine. Streams trusted postcode/city/street rows from
   `preload.csv` into `targets` (no `uniform` normalization).
2. **`preloadFeed`** (`preload.go`) — started after the trusted preload
   finishes. Streams validated postcode/city pairs from the CSV into
   `targets`.
3. **`collect`** (`collect.go`) — spun up in its own goroutine before the\n   feeds; ranges over `targets` until
   the channel is closed, building the in-memory place index and emitting
   JSON. Single goroutine; all dedup/correction logic lives here.
4. **`pbfparser`** (`parser.go`) — started only **after** `preload.Wait()`
   returns. Decodes the PBF and pushes qualifying tag sets into `targets`.
5. After the parser finishes, `close(targets)` lets `collect` drain and exit.

Ordering matters: preload must finish before the parser starts because both
feed the same unbuffered `targets` channel and `collect` is the sole
consumer. The channel is unbuffered; throughput depends on `collect`
keeping up.

### Concurrency

- `core.go` declares **package-level** `sync.WaitGroup`s (`preload`,
  `parser`, `collector`) and the `targets` channel. They are global state —
  `Parse` is not safe to call concurrently or more than once per process.
  Stages are launched via `sync.WaitGroup.Go` (Go 1.25+); the stage functions
  no longer call `Done()` themselves.
- `Target.Worker` exists in the struct but is **unused** for parser
  parallelism; the PBF decoder parallelism is controlled separately inside
  `internal/pbf` via `runtime.GOMAXPROCS` (`DefaultNCpu` = max(GOMAXPROCS-1, 1)).
- The PBF decoder uses `github.com/destel/rill` for pipelined
  blob→batch→object processing; the `rill.Discard` call in
  `internal/pbf/decoder.go` drains the background pipeline on `Decoder.Close`.

### PBF decoding stack (`internal/`)

- `internal/decoder` — low-level PBF blob/header parsing; converts
  `PrimitiveBlock`s into `model.Object`s. Note: only **Nodes** and
  **DenseNodes** are decoded into elements; `parser.go` switches on
  `*model.Node`, `*model.Way`, `*model.Relation`, `*model.Header` (Way,
  Relation and Header cases are no-ops — the tool effectively only processes
  nodes with `addr:*` tags); the `default` case panics on truly unknown types.
- `internal/pbf` — public `Decoder` over `internal/decoder`, configurable via
  functional options (`WithProtoBufferSize`, `WithProtoBatchSize`,
  `WithNCpus`). Default batch size 16, buffer 1 MiB.
- `internal/protobuf` — generated from `internal/protobuf/osm.proto`
  (`proto2`, `option go_package = "m4o.io/pbf/protobuf"` — note the
  go_package path is **stale** and does not match the actual import path).
  Regenerate with `protoc` if the schema changes; the generated
  `osm.pb.go` is committed.
- `internal/model` — shared `Node`/`Way`/`Relation`/`Header` types and the
  `Object` sealed-via-`foo()` interface. `ElementType` has a
  `//go:generate stringer` directive (`elementtype_string.go` is committed).
- `internal/core` — `sync.Pool`-backed `PooledBuffer` for byte reuse.

## Address Normalization (`uniform`)

`uniform.go` dispatches by `t.Country`; only `"DE"` is implemented
(`uniform_de.go`). Adding a new country means adding a `case` in
`uniform.go` and a `uniform<CC>.go` file implementing the equivalent of
`uniformDE`. The returned `int` is a **count of normalization failures**,
not a success flag — it is accumulated into `OSM:PBF:Err:Uniform`.

DE normalization:
- Rejects non-Latin1 city/street (reported and counted into
  `OSM:PBF:Err:Uniform`).
- Parses postcode as int; pads 4-digit postcodes with a leading `0`,
  rejects anything else as failure.
- `tryNormStreetDE` canonicalizes `Strasse`/`Str.`/`strasse`/`str.` →
  `Straße`/`straße`. Order of replacements is significant.
- `tryNormCityDE` expands abbreviations (`a.d.` → `an der`, `i.` → `in`,
  `b.` → `bei`, etc.), fixes common typos (`in Isartal` → `im Isartal`,
  `/ Saale` → `/Saale`), then `camelCaseCityDE` title-cases words while
  preserving a hardcoded list of lowercase particles (`an`, `am`, `der`,
  `in`, `von`, …) and uppercasing `ii`→`II`, `ot`→`OT`, `nb`→`NB`, etc.
  Particles are matched on the **lowercased** token, so any new exception
  must be added lowercase.

## Place ID & Dedup

- `id(in)` (`helper.go`) = first 12 bytes of `sha256(in + _secretMAC)` where
  `_secretMAC` is a hardcoded constant. The place ID is the
  `country+postcode+city+street` concatenation. **IDs are not stable across
  versions if the secret changes**; the secret is committed in the clear —
  do not treat it as confidential, but be aware changing it invalidates all
  previously generated mapping tables.
- `collect.go` dedups by `placeID` first, then by nested
  `postcode→city→street` maps. It also performs auto-correction for cities
  that differ only by separator (` `, `-`) against an existing entry, and
  emits Levenshtein-distance<2 warnings to `warning.json` (not corrections).
- Three parallel index structures are maintained in `collect`:
  `placeIDs` (placeID→tagSet), `places` (postcode→city→street→placeIdHex),
  `places2` (postcode→city→street→bool, written as `addr.json`).

## Output Files

All written by `writeJsonFile` (`writer.go`) into `json/<COUNTRY_CODE>/`,
created with `MkdirAll(0755)`, files `0644`, `json.MarshalIndent` with `\t`.
Files are overwritten silently:

| File | Source | Content |
| --- | --- | --- |
| `addr.json` | `places2` | postcode→city→street→true |
| `addr2placeID.json` | `places` | postcode→city→street→hex placeID |
| `placeID2addr.json` | `p` | hexPlaceID→tagSet |
| `warning.json` | `warning` | message→count (Levenshtein + non-preloaded postcode/city) |
| `corrected.json` | `corrected` | message→count (auto-corrected city separators) |
| `error.preload.json` | `fail` (preload only) | message→count (CSV length + CSV parse errors) |

`examples/json/DE/` contains sample outputs for reference; do not edit by
hand — they are regenerated by `make -C examples`.

## Conventions & Gotchas

- **Error handling style**: panics on internal/decode/JSON-marshal errors;
  `log.Fatal` on CLI/setup errors in `main.go`; per-record validation
  failures are counted and routed to a `*.json` error file instead of
  aborting. Match this pattern: don't return errors from per-record logic.
- **Logging**: progress is emitted via `fmt.Printf` with leading `\n` and
  the `OSM:<Section>:<Key>` tag format padded to fixed column widths
  (see existing `OSM:PBF:Parsed:Objects    # ...` lines). `hu()` in
  `helper.go` right-pads integer stats to 11 chars using German locale
  separators. Preserve this formatting when adding new stat lines.
- **Receiver names**: short, lowercase, single-word or single-letter
  (`t *tagSet`, `r Node`, `b *PooledBuffer`, `d *Decoder`). Match the
  surrounding file.
- **Exported symbols** in the public package (`Target`, `Parse`) are the
  only intended API surface; the `internal/` tree is intentionally
  unexported. `cmd/osm2addr/main.go` is the reference consumer.
- **`Target.PreLoad` field indices** are hardcoded in `checkPreloadFile`
  for the German CSV schema. Supporting a different preload schema requires
  parameterizing `Fields`/`City`/`Postcode`/`PostcodeLenght` — note the
  typo `Lenght` (sic) in the field name; do not "fix" it without renaming
  every usage.
- **Debug logging**: set the `DEBUG` environment variable to any non-empty
  value (except `0`, `false`, `no`, `off`) to enable per-event
  `OSM:DEBUG:*` lines (`debug.go`).
- **`data/` is fully gitignored** (`data/.gitignore` is `*` except
  `.gitignore` itself). PBF files are large; never commit them. Output
  `json/` directories under `cmd/osm2addr/` and `examples/` are also
  gitignored except committed samples under `examples/json/`.
- **`make deps` is destructive** — it deletes `go.mod` and `go.sum` and
  re-runs `go mod init`/`tidy`. Do not run it casually; dependabot manages
  dependency bumps via PRs (see `.github/dependabot.yml`).
- **CI Go version**: `.github/workflows/golang.yml` uses `go-version: stable`
  (aligned with `go.mod`, which tracks the current Go release line). The
  golangci-lint workflow also uses `stable`.
- **No comments policy in this repo**: most files have no inline comments;
  doc comments are terse (`// uniform ...`). When editing, follow the
  existing minimal style rather than adding explanatory commentary.
- **`nolint` directives** exist (e.g. `//nolint:forcetypeassert` in
  `internal/core/buffer.go`); preserve them when refactoring.

## Memory / Project-Specific Notes

- The `secret` in `_secretMAC` is not configurable via CLI; changing it
  requires a recompile and invalidates downstream mapping tables.
- The protobuf `go_package` option in `osm.proto` says `m4o.io/pbf/protobuf`
  — this is a leftover from upstream and does **not** match the actual
  import path `paepcke.de/osm2addr/internal/protobuf`. Regenerating the
  `.pb.go` without fixing this option can break the build.
- `examples/json/DE/` and `examples/README.md` are regenerated by
  `make -C examples`; if their prose ever lags the code, treat the code as
  the source of truth.
