<div align="center">

# 🗺️ osm2addr

**One PBF in. Clean postal-address tables out.**

Turn an OpenStreetMap extract into deduplicated, normalized
`postcode → city → street` mapping tables — deterministic place IDs,
auto-corrected city names, and a warning feed for everything that does not
look like official postal data.

[![Go Reference](https://pkg.go.dev/badge/paepcke.de/osm2addr.svg)](https://pkg.go.dev/paepcke.de/osm2addr)
[![Go Report Card](https://goreportcard.com/badge/paepcke.de/osm2addr)](https://goreportcard.com/report/paepcke.de/osm2addr)
[![Go Build](https://github.com/paepckehh/osm2addr/actions/workflows/golang.yml/badge.svg)](https://github.com/paepckehh/osm2addr/actions/workflows/golang.yml)
[![golangci-lint](https://github.com/paepckehh/osm2addr/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/paepckehh/osm2addr/actions/workflows/golangci-lint.yml)
[![License](https://img.shields.io/github/license/paepckehh/osm2addr)](https://github.com/paepckehh/osm2addr/blob/main/LICENSE)
[![SemVer](https://img.shields.io/github/v/release/paepckehh/osm2addr)](https://github.com/paepckehh/osm2addr/releases/latest)
[![built with nix](https://builtwithnix.org/badge.svg)](https://search.nixos.org/packages?channel=unstable&from=0&size=50&sort=relevance&type=packages&query=osm2addr)

`go run paepcke.de/osm2addr/cmd/osm2addr@latest`

</div>

---

## 💡 Why

Address validation databases want the boring part of OpenStreetMap: *this
postcode belongs to this city, and this is the list of streets in it*. They do
not want 400 million nodes, duplicate nodes with the same address, or
`Hauptstr.` next to `Hauptstraße`.

`osm2addr` reads a `.osm.pbf` extract, keeps only complete `addr:*` nodes,
normalizes them per country, deduplicates by a deterministic place ID and emits
small JSON tables you can import directly into MongoDB, PostgreSQL, SQLite or
whatever your validator uses.

Measured on a consumer laptop: roughly **400 million nodes in under 50
seconds**, single binary, no network access.

## ✨ Features

- ⚡ **Fast** — pipelined PBF decoding (`rill`), blob → batch → object, parallel
  by default.
- 🧼 **Normalized** — country-aware canonicalization of streets, cities,
  abbreviations and postcodes. German (`DE`) is implemented; adding a country
  is one dispatch case and one file.
- 🎯 **Deterministic place IDs** — 96-bit SHA-256 prefix of the normalized
  place, so the same address always maps to the same ID.
- 🧠 **Auto-correction** — `Bad Homburg` and `Bad-Homburg` collapse into one
  city; near-duplicate city spellings are flagged by Levenshtein distance.
- 🧾 **Warnings, not crashes** — non-preloaded postcodes, preload length
  errors, bad CSV rows and normalization drops are counted into JSON error
  files. The run never aborts on a single bad record.
- 📥 **Two preload feeds** — an optional validated postcode/city CSV, plus a
  trusted street-level `preload.csv` (header auto-detected) that seeds the
  index before OSM data.
- 🔌 **App *and* library** — the CLI is a thin wrapper around
  `osm2addr.Parse`; every stage is testable in isolation.
- 🪶 **Pure Go, tiny deps** — `CGO_ENABLED=0`, static binaries for
  linux/freebsd/darwin/netbsd/openbsd/windows, amd64 and arm64.
- 🧪 **Fast tests** — parallel, self-contained, with a tiny generated PBF
  fixture. `go test ./...` finishes in milliseconds.
- 🐧 **Grep-friendly** — stable `OSM:<Section>:<Key>` progress lines, optional
  per-record debug output.

Sibling projects with the same look and feel:
[`certinfo`](https://paepcke.de/certinfo/), [`tlsinfo`](https://paepcke.de/tlsinfo/),
[`dnsinfo`](https://paepcke.de/dnsinfo/).

## 🚀 Install

**Go** (any platform with a toolchain):

```shell
go install paepcke.de/osm2addr/cmd/osm2addr@latest
```

**No toolchain? Run it straight away:**

```shell
go run paepcke.de/osm2addr/cmd/osm2addr@latest DE data/germany-latest.osm.pbf
```

**Pre-built binaries** (linux/freebsd/darwin/netbsd/openbsd/windows,
amd64+arm64):

→ [github.com/paepckehh/osm2addr/releases/latest](https://github.com/paepckehh/osm2addr/releases/latest)

**NixOS / nix:**

```shell
nix-shell -p osm2addr
```

```nix
environment.systemPackages = [ pkgs.osm2addr ];
```

## 🎬 Quickstart

```shell
mkdir -p data
curl --location --output data/germany-latest.osm.pbf \
  https://download.geofabrik.de/europe/germany-latest.osm.pbf

osm2addr DE data/germany-latest.osm.pbf
```

Output lands in `./json/DE/`:

```text
json/DE/addr.json                  postcode → city → street → true
json/DE/addr2placeID.json          postcode → city → street → place ID
json/DE/placeID2addr.json          place ID → {country, postcode, city, street}
json/DE/warning.json               near-duplicates and non-preloaded pairs
json/DE/corrected.json             auto-corrected city separators
json/DE/error.preload.json         validated preload row errors
json/DE/error.preload.trusted.json trusted preload row errors
```

The last two files are only written when the matching preload file exists.

Or use the Makefile shortcuts:

```shell
make update-de     # download the Germany extract into data/
make build         # build ./osm2addr
make run           # build and run with defaults
make test          # race-enabled test suite
```

## 🎛 Usage

```text
osm2addr [<COUNTRY_CODE> [<PATH_TO_PBF>]]
```

| Argument | Description | Default |
| --- | --- | --- |
| `COUNTRY_CODE` | Exactly 2 characters, uppercased internally (e.g. `DE`) | `DE` |
| `PATH_TO_PBF` | OSM PBF extract, relative to the process CWD | `data/germany-latest.osm.pbf` |

Output is always written to `json/<COUNTRY_CODE>/` relative to the process
CWD, never next to the binary.

The PBF file is opened exactly once after argument parsing. If it cannot be
read, the process exits with a fatal error before any output is written.

### 🐞 Debug logging

Set `DEBUG` to any non-empty value except `0`, `false`, `no`, `off`:

```shell
DEBUG=1 osm2addr DE data/germany-latest.osm.pbf
```

Every dropped, corrected or added record becomes an `OSM:DEBUG:*` line.
`DEBUG` only adds logging; it never changes what is written.

### 📊 Sample run

```text
OSM:Version              # osm2addr v0.0.37 (commit deadbeef, built 2026-10-24T00:00:00Z)
OSM:Startup               # 2026-10-24 09:23:45 +0000 UTC m=+0.000386754
OSM:TargetCountry         # DE
OSM:Debug                 # false
OSM:File                  # data/germany-latest.osm.pbf
----------------------------------------------------------------------------------
OSM:PreLoadTrusted        # data/preload.csv
OSM:PreLoadTrusted:Fail   #           0
OSM:PreLoadTrusted:Total  #     119.733
OSM:Writer:JSON           # json/DE/error.preload.trusted.json
----------------------------------------------------------------------------------
OSM:PreLoadFile           # data/validated-preload/DE.csv
OSM:PreLoadFile:Fail      #           1
OSM:PreLoadFile:Total     #      12.853
OSM:Writer:JSON           # json/DE/error.preload.json
----------------------------------------------------------------------------------
OSM:PBF:File:URL          # https://download.geofabrik.de/europe/germany-updates
OSM:PBF:File:Repl:USM     # 4330
OSM:PBF:File:Repl:TS      # 2026-10-20 21:21:14 +0000 UTC
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

| Counter | Meaning |
| --- | --- |
| `PBF:Parsed:Objects` | Nodes, ways, relations and headers seen |
| `PBF:Parsed:Tags` | Tags inspected, including non-address tags |
| `PBF:Complete:AddrTags` | Nodes with complete postcode + city + street |
| `PBF:Err:Uniform` | Entries dropped by country normalization |
| `PBF:Err:Country/Postcode/City/Street` | Malformed address fields |
| `Corrected:Auto:*` | City separators normalized against an existing city |
| `Corrected:Warn:*` | Levenshtein near-duplicates and non-preloaded pairs |
| `Collect:Places:Total` | Unique place IDs exported |

Nodes without an `addr:country` tag inherit the target country. A node that
*has* an `addr:country` tag but whose value is malformed (not exactly two
uppercase characters) is dropped instead of falling back, so a lowercase
`at` cannot be silently imported as a German address.

## 📥 Preload data

Preloads are optional. They pre-seed the place index so that OSM data can be
checked against a known-good address list. Both feeds run to completion before
the PBF parser starts.

### Validated preload CSV (expected for `DE`)

Place the file at `<dirname of PBF>/validated-preload/<COUNTRY_CODE>.csv`:

- comma separated, **6 fields per record**, header row skipped;
- city in column **2**, postcode in column **3**;
- postcode must be 5 characters (German).

```csv
osm_id,ags,ort,plz,landkreis,bundesland
1104550,08335001,Aach,78267,Landkreis Konstanz,Baden-Württemberg
```

Rows with the wrong field count, bad quoting or failing length checks are
counted into `json/<CC>/error.preload.json`. Their entries take part in dedup
and place-ID generation like OSM data.

A row may omit the street. Such a city-only row seeds the postcode/city pair
but produces no place ID and no street entry.

A path that opens but is not a regular file is rejected up front, because
`csv.Reader` would otherwise fail on every read without ever reaching
`io.EOF`.

### Trusted preload.csv

A file named `preload.csv` next to the PBF (process CWD as fallback) is loaded
**before** everything else:

- the first line is the header; the separator is detected (`,` default,
  `;`, tab and `|` supported), a UTF-8 BOM is stripped and UTF-8 is verified;
- the header must contain the exact names `POSTLEITZAHL`, `ORT_NAME` and
  `STRASSE_NAME`, in any order;
- rows are trusted street-level values: **no** normalization and no matching,
  but stray surrounding whitespace is trimmed off every value;
- rows with an empty postcode or city, or broken CSV, are counted into
  `json/<CC>/error.preload.trusted.json`.

```csv
POSTLEITZAHL;ORT_NAME;STRASSE_NAME
24937;Flensburg;Am Pferdewasser
```

A production-shaped sample lives in
[example-preload.csv](example-preload.csv).

## 📤 Output files

All files are written into `json/<COUNTRY_CODE>/`, tab-indented, `0644`,
overwritten silently. Directories are created with `0755`.

| File | Content |
| --- | --- |
| `addr.json` | `postcode → city → street → true` |
| `addr2placeID.json` | `postcode → city → street → hex place ID` |
| `placeID2addr.json` | `hex place ID → {country, postcode, city, street}` |
| `warning.json` | message → count (Levenshtein near-duplicates, non-preloaded postcode/city) |
| `corrected.json` | message → count (auto-corrected city separators) |
| `error.preload.json` | message → count (validated preload row errors) |
| `error.preload.trusted.json` | message → count (trusted preload row errors) |

```json
{
  "01067": {
    "Dresden": {
      "Adlergasse": "05ba9df31e1235fbd3d23296"
    }
  }
}
```

`addr.json` keeps the same shape with `true` instead of the place ID.
`warning.json` and `corrected.json` are always written, even when empty (`{}`).
Sample outputs for a German extract live in [examples/json/DE](examples/json/DE).

### 🆔 Place ID

A place ID is the first 12 bytes of
`sha256(country \x00 postcode \x00 city \x00 street + secretMAC)`,
hex-encoded to 24 characters. IDs are deterministic for the same normalized
address.

> **Breaking:** the NUL-separated key format was introduced in `v0.0.36` to
> make city/street boundary collisions impossible. Regenerate previously
> imported tables. IDs also change whenever `_secretMAC` changes, which is why
> the secret is a compile-time constant.

## 🏗 Architecture

```text
preload.csv (trusted)  ──┐
validated preload (csv) ──┼──>  chan *tagSet  targets  ──>  collect  ──> JSON files
pbfparser (pbf)        ───┘
```

1. **`trustedPreloadFeed`** streams trusted postcode/city/street rows first.
2. **`preloadFeed`** streams validated postcode/city pairs next.
3. **`pbfparser`** decodes the PBF and pushes complete `addr:*` tag sets.
4. **`collect`** is the only consumer: dedup, city auto-correction,
   Levenshtein warnings and JSON emission.

Both preloads finish before the parser starts. `Parse` creates the channel
locally, so the package holds **no mutable global state** and each stage runs
independently in tests.

The `internal/` tree handles PBF decoding (`pbf` → `decoder` → `protobuf`),
with `github.com/destel/rill` for the blob → batch → object pipeline and
parallelism from `runtime.GOMAXPROCS`. Only nodes and dense nodes are decoded;
ways and relations are ignored.

```
internal/pbf       public Decoder, functional options, rill pipeline
internal/decoder   blob/header parsing, PrimitiveBlock → model.Object
internal/protobuf  generated osm.proto types (committed)
internal/model     Node/Way/Relation/Header and ElementType
internal/core      pooled byte buffers
```

### 🔤 Address normalization (`DE`)

`uniform.go` dispatches by country; `uniform_de.go` implements:

- reject non-Latin1 city/street;
- parse the postcode as int, pad 4-digit codes with `0`, reject anything
  outside `0`–`99999`;
- canonicalize `Strasse`/`Str.`/`strasse`/`str.` → `Straße`/`straße`
  (`Hauptstr.` → `Hauptstraße`, no phantom space);
- expand abbreviations (`a.d.` → `an der`, `i.` → `in`, `b.` → `bei`, …), fix
  common typos (`in Isartal` → `im Isartal`, `/ Saale` → `/Saale`), collapse
  whitespace, then title-case while preserving lowercase particles (`an`,
  `am`, `der`, `in`, `von`, …) and uppercasing `ii`, `ot`, `nb`, `ofr.`,
  `opf.`.

Adding a country: add a `case` to `uniform.go` and a `uniform<CC>.go` file
returning the drop flag.

## 📦 Use it as a library

```go
package main

import (
	"log"
	"os"

	"paepcke.de/osm2addr"
)

func main() {
	f, err := os.Open("data/germany-latest.osm.pbf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	osm2addr.InitDebug()
	if err := osm2addr.Parse(&osm2addr.Target{
		Country:  "DE",
		File:     f,
		FileName: "data/germany-latest.osm.pbf",
	}); err != nil {
		log.Fatal(err)
	}
}
```

Public API ([godoc](https://pkg.go.dev/paepcke.de/osm2addr)):

- `Parse(*Target) error` — run the full pipeline.
- `Target` — country, open PBF file plus its filename (used to locate the
  preload files).
- `InitDebug()` / `DebugEnabled()` — control the `OSM:DEBUG:*` event log.

`Parse` writes the JSON tables relative to the working directory, exactly like
the CLI. The `internal/` tree is not part of the public API.

## 🛠 Development

```shell
make build       # go build with version/commit/buildtime ldflags
make check       # gofmt + go vet + go fix + staticcheck ./...
make test        # CGO_ENABLED=1 go test -race -count=1 ./...
make update-de   # download the Germany extract into data/
make update-dach # download the DACH extract into data/
make update-eu   # download the Europe extract into data/
make deps        # destructive: recreates go.mod/go.sum (avoid; dependabot owns deps)
```

Tests are fast and parallel, and never touch `data/` or the repo `json/` tree
(pure fixtures, `t.TempDir()`):

```shell
go test -count=1 ./...        # unit + pipeline tests, tiny generated PBF
go test -race -count=1 ./...
```

Every change follows the fixed flow:

1. `gofmt -w .`
2. `go build -o osm2addr ./cmd/osm2addr`
3. `go test -count=1 ./...`
4. commit
5. tag `v0.0.<N+1>` (patch only; never move or reuse a tag)

Repository layout:

```text
cmd/osm2addr/   CLI entrypoint
*.go            pipeline, normalization, collector, writer
internal/       PBF decoding stack (Apache-2.0)
examples/       sample output and regeneration Makefile
data/           gitignored PBF/extract workspace
```

## 📚 Docs

- API reference: [pkg.go.dev/paepcke.de/osm2addr](https://pkg.go.dev/paepcke.de/osm2addr)
- Project home: [paepcke.de/osm2addr](https://paepcke.de/osm2addr/)
- Sample output and run instructions: [examples/](examples/)
- Preload example: [example-preload.csv](example-preload.csv)
- Releases: [github.com/paepckehh/osm2addr/releases](https://github.com/paepckehh/osm2addr/releases)

## 🤝 Contributing

Yes, please! Open an issue or send a PR. Keep tests parallel and
self-contained, and run the fixed flow above before submitting. See
[AGENTS.md](AGENTS.md) for the repository conventions.

## 🧭 Roadmap

- [ ] Additional country normalizers beside `DE`
- [ ] Configurable preload schema via `Target` fields
- [ ] Optional per-country street-level validation reports

## 🛡 License

[BSD 3-Clause](LICENSE). The `internal/` PBF decoding stack is Apache-2.0 —
see [internal/LICENSE](internal/LICENSE).

## 📃 Citation

```bibtex
@misc{osm2addr,
  author = {Michael Paepcke},
  title = {osm2addr: turn OpenStreetMap PBF extracts into normalized postal-address mapping tables},
  year = {2024},
  publisher = {GitHub},
  journal = {GitHub repository},
  howpublished = {\url{https://paepcke.de/osm2addr}}
}
```
