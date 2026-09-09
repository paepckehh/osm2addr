# osm2addr

[![Go Reference](https://pkg.go.dev/badge/paepcke.de/osm2addr.svg)](https://pkg.go.dev/paepcke.de/osm2addr)
[![Go Report Card](https://goreportcard.com/badge/paepcke.de/osm2addr)](https://goreportcard.com/report/paepcke.de/osm2addr)
[![Go Build](https://github.com/paepckehh/osm2addr/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/osm2addr/actions/workflows/golang.yml)
[![golangci-lint](https://github.com/paepckehh/osm2addr/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/paepckehh/osm2addr/actions/workflows/golangci-lint.yml)
[![License](https://img.shields.io/github/license/paepckehh/osm2addr)](https://github.com/paepckehh/osm2addr/blob/master/LICENSE)
[![SemVer](https://img.shields.io/github/v/release/paepckehh/osm2addr)](https://github.com/paepckehh/osm2addr/releases/latest)
[![built with nix](https://builtwithnix.org/badge.svg)](https://search.nixos.org/packages?channel=unstable&from=0&size=50&sort=relevance&type=packages&query=osm2addr)

`osm2addr` parses [OpenStreetMap](https://www.openstreetmap.org/) PBF extract
files and emits normalized, postal-address mapping tables (JSON) keyed by a
deterministic place ID. It is optimized for the German (`DE`) address dataset,
but the architecture is country-aware via the `uniform` dispatch and can be
extended by adding a per-country normalizer. The output is intended to feed
address validation databases (MongoDB, etc.).

## Table of contents

- [Features](#features)
- [Install](#install)
- [Quickstart](#quickstart)
- [Usage](#usage)
- [Preload data](#preload-data)
- [Output files](#output-files)
- [Architecture](#architecture)
- [Development](#development)
- [License](#license)

## Features

- **Fast** — extracts, analyzes and normalizes around 400.000.000 OSM nodes
  in under 50 sec on a consumer laptop.
- **Deduplicated** — every place (postcode + city + street) is stored exactly
  once, keyed by a deterministic SHA-256-derived place ID.
- **Normalized** — country-specific canonicalization of street names, city
  names, abbreviations and postcodes (German `DE` implemented; extensible).
- **Auto-correction** — cities that differ only by separator choice
  (`Bad Homburg` vs `Bad-Homburg`) are collapsed automatically; near-duplicate
  city names within the same postcode (Levenshtein distance < 2) are flagged.
- **Preload support** — optional validated and trusted CSV inputs pre-seed the
  place index so non-preloaded entries are surfaced as warnings.
- **Offline** — no dependency on external services; runs fully offline against
  a local `.osm.pbf` file.
- **Tested** — parallel, self-contained unit tests with tiny in-memory
  fixtures (including a minimal generated PBF file), runnable with
  `go test ./...` in milliseconds.

## Install

```
go install paepcke.de/osm2addr/cmd/osm2addr@latest
```

Or build from source:

```
git clone https://github.com/paepckehh/osm2addr
cd osm2addr
make build
```

## Quickstart

```
mkdir -p data
curl --output data/germany-latest.osm.pbf https://download.geofabrik.de/europe/germany-latest.osm.pbf
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE data/germany-latest.osm.pbf
```

## Usage

```
osm2addr [<COUNTRY_CODE> [<PATH_TO_PBF>]]
```

| Argument | Description | Default |
| --- | --- | --- |
| `COUNTRY_CODE` | Exactly 2 characters, uppercased internally (example: `DE`) | `DE` |
| `PATH_TO_PBF` | Path to an OSM PBF extract, relative to the process CWD | `data/germany-latest.osm.pbf` |

### Verbose logging

Set the `DEBUG` environment variable to any non-empty value (except `0`,
`false`, `no`, `off`) to enable per-event debug logging:

```
DEBUG=1 osm2addr DE data/germany-latest.osm.pbf
```

### Sample output

```
OSM:Version              # osm2addr dev (commit none, built unknown)
OSM:Startup               # 2025-03-10 09:23:45.921159027 +0000 UTC m=+0.000386754
OSM:TargetCountry         # DE
OSM:Debug                 # false
OSM:File                  # data/germany-latest.osm.pbf
----------------------------------------------------------------------------------
OSM:PreLoadFile           # data/validated-preload/DE.csv
OSM:PreLoadFile:Fail      # 1
OSM:PreLoadFile:Total     # 12853
OSM:Writer:JSON           # json/DE/error.preload.json
----------------------------------------------------------------------------------
OSM:PBF:File:URL          # https://download.geofabrik.de/europe/germany-updates
OSM:PBF:File:Repl:USM     # 4330
OSM:PBF:File:Repl:TS      # 2025-02-13 21:21:14 +0000 UTC
----------------------------------------------------------------------------------
OSM:PBF:Parsed:Objects    # 411.113.874
OSM:PBF:Parsed:Tags       #  71.720.811
OSM:PBF:Parsed:Country    #   2.792.983
OSM:PBF:Parsed:Street     #   3.905.365
OSM:PBF:Parsed:City       #   3.659.945
OSM:PBF:Parsed:Postcode   #   3.513.835
----------------------------------------------------------------------------------
OSM:PBF:Complete:AddrTags #   2.711.889
OSM:PBF:Uniq:Country      #          10
OSM:PBF:Err:Uniform       #      21.854
OSM:PBF:Err:Country       #           0
OSM:PBF:Err:Postcode      #           0
OSM:PBF:Err:City          #          17
OSM:PBF:Err:Street        #         162
----------------------------------------------------------------------------------
OSM:Corrected:Auto:Cases  #          78
OSM:Corrected:Auto:Total  #      21.082
OSM:Corrected:Warn:Cases  #           8
OSM:Corrected:Warn:Total  #           8
OSM:Collect:Places:Total  #     278.381
----------------------------------------------------------------------------------
OSM:Writer:JSON           # json/DE/addr.json
OSM:Writer:JSON           # json/DE/addr2placeID.json
OSM:Writer:JSON           # json/DE/placeID2addr.json
OSM:Writer:JSON           # json/DE/warning.json
OSM:Writer:JSON           # json/DE/corrected.json
OSM:Time:Total            # 45.329429572s
```

## Preload data

### Validated preload CSV (optional, expected for `DE`)

If a file exists at `<dirname of PBF>/validated-preload/<COUNTRY_CODE>.csv`, it
is parsed as a comma-separated CSV with **6 fields per record**, a header row
(skipped), and fixed column indices:

| Field    | Column |
| ---      | ---    |
| City     | 2      |
| Postcode | 3      |

Expected postcode length is **5** (German). Records failing length checks —
and CSV rows that fail to parse (e.g. wrong field count, stray quote) — are
counted into `json/<CC>/error.preload.json` instead of aborting the run. The
preload entries participate in dedup and place-ID generation like any other
input.

### Trusted preload CSV (optional)

If a file named `preload.csv` exists at program start (next to the PBF file,
or in the process working directory as fallback), it is loaded as trusted
start values before any other input:

- The first line is assumed to be a valid CSV header; the separator is
  detected from the header line (`,` is the assumed default; `;`, tab and `|`
  are detected), a UTF-8 BOM is stripped and UTF-8 encoding is confirmed.
- The header must contain the exact field names `POSTLEITZAHL`, `ORT_NAME`
  and `STRASSE_NAME`. If any of them is missing, the file is rejected.
- Its rows are fed into the pipeline **first** and are treated as trusted
  postcode/city/street values: no country-specific normalization or matching
  is applied. They pre-seed all place index arrays before the OSM import data
  is processed in the usual way.
- Rows that fail CSV parsing are counted into
  `json/<CC>/error.preload.trusted.json` instead of aborting the run.

## Output files

All JSON files are written to `json/<COUNTRY_CODE>/` (relative to the process
CWD) with tab indentation and `0644` permissions; existing files are
overwritten silently.

| File | Content |
| --- | --- |
| `addr.json` | `postcode` → `city` → `street` → `true` |
| `addr2placeID.json` | `postcode` → `city` → `street` → hex place ID |
| `placeID2addr.json` | hex place ID → `{country, postcode, city, street}` |
| `warning.json` | message → count (Levenshtein near-duplicates + non-preloaded postcode/city) |
| `corrected.json` | message → count (auto-corrected city separators) |
| `error.preload.json` | message → count (preload CSV length errors) |
| `error.preload.trusted.json` | message → count (trusted preload.csv CSV parse errors) |

### Place ID

A place ID is the first 12 bytes of `sha256(country + postcode + city + street
+ secretMAC)`, hex-encoded. IDs are deterministic but **not stable across
versions if the secret changes**; the secret is compiled in and changing it
invalidates all previously generated mapping tables.

## Architecture

The pipeline is a fixed producer/consumer chain wired in `core.go::Parse`:

```
preload.csv (trusted)  ──┐
validated preload (csv) ──┼──>  chan *tagSet  targets  ──>  collect ──> JSON files
pbfparser (pbf)        ───┘
```

1. **`trustedPreloadFeed`** streams trusted postcode/city/street rows from
   `preload.csv` into the `targets` channel, ahead of all other input.
2. **`preloadFeed`** streams validated postcode/city pairs from the validated
   preload CSV into `targets`.
3. **`pbfparser`** decodes the PBF and pushes qualifying `addr:*` tag sets
   into `targets`.
4. **`collect`** drains `targets`, building the in-memory place index and
   emitting JSON. All dedup and auto-correction logic lives here.

Both preload feeds run to completion before the parser starts; all producers
share one unbuffered `targets` channel and `collect` is its sole consumer.
`Parse` creates the channel locally, so the package keeps no mutable global
state and every stage is independently testable.

PBF decoding is handled by the `internal/` tree (`pbf` → `decoder` →
`protobuf`), which uses `github.com/destel/rill` for pipelined
blob → batch → object processing with parallelism controlled via
`runtime.GOMAXPROCS`.

### Address normalization

`uniform.go` dispatches by `t.Country`; only `"DE"` is implemented
(`uniform_de.go`). Adding a new country means adding a `case` in `uniform.go`
and a `uniform<CC>.go` file. DE normalization:

- Rejects non-Latin1 city/street.
- Parses the postcode as int; pads 4-digit postcodes with a leading `0`,
  rejects anything outside `0`-`99999`.
- Canonicalizes `Strasse`/`Str.`/`strasse`/`str.` → `Straße`/`straße`.
- Expands abbreviations (`a.d.` → `an der`, `i.` → `in`, `b.` → `bei`, …),
  fixes common typos, then title-cases words while preserving a hardcoded list
  of lowercase particles.

## Development

```
make build     # build the binary
make check     # gofmt + go vet + go fix + staticcheck
make deps      # destructive: recreates go.mod/go.sum (avoid; dependabot manages deps)
make update-de # download the Germany PBF extract
```

Test suite (fast, parallel, no large fixtures required):

```
go test -count=1 ./...        # unit + pipeline tests, incl. a tiny generated PBF
go test -race -count=1 ./...  # with the race detector (requires CGO_ENABLED=1)
```

The fixed workflow for every change:

1. `gofmt -w .`
2. `go build -o osm2addr ./cmd/osm2addr`
3. `go test -count=1 ./...`
4. commit
5. tag `v0.0.<N+1>` (patch segment only; never move or reuse a tag)

## License

BSD 3-Clause — see [LICENSE](LICENSE). The `internal/` PBF decoding stack is
Apache-2.0 — see [internal/LICENSE](internal/LICENSE).
