# 📚 Examples

Sample outputs generated from the German OpenStreetMap extract. Use these to
inspect the JSON shape before importing the tables into your own database.

```text
json/DE/addr2placeID.json   postcode → city → street → place ID
json/DE/placeID2addr.json   place ID → {country, postcode, city, street}
json/DE/warning.json        Levenshtein near-duplicates and non-preloaded pairs
json/DE/corrected.json      auto-corrected city separators
json/DE/error.preload.json  validated preload row errors
```

`addr.json` (`postcode → city → street → true`) is generated on every run but
is not committed here because it duplicates the key structure of
`addr2placeID.json`.

## 🔁 Regenerate

From the repository root:

```shell
make update-de                  # download data/germany-latest.osm.pbf
make -C examples                # regenerate examples/json/DE/
```

Or manually:

```shell
mkdir -p data/validated-preload
curl --output data/germany-latest.osm.pbf \
  https://download.geofabrik.de/europe/germany-latest.osm.pbf
curl --output data/validated-preload/DE.csv \
  https://downloads.suche-postleitzahl.org/v2/public/zuordnung_plz_ort.csv
go run ../cmd/osm2addr/main.go DE ../data/germany-latest.osm.pbf
```

The run writes into `examples/json/DE/` and prints the usual
`OSM:<Section>:<Key>` progress report. A full German extract takes roughly 45
seconds on a consumer laptop.

## 🧾 Shape

`addr2placeID.json`:

```json
{
  "01067": {
    "Dresden": {
      "Adlergasse": "05ba9df31e1235fbd3d23296"
    }
  }
}
```

`placeID2addr.json`:

```json
{
  "05ba9df31e1235fbd3d23296": {
    "country": "DE",
    "postcode": "01067",
    "city": "Dresden",
    "street": "Adlergasse"
  }
}
```

`warning.json` and `corrected.json` are message → count maps:

```json
{
  "[WARNING][LEVENSHTEIN:1][POSTCODE:04779][CITY]#Wermsdorf#Wermsorf#": 1
}
```

Place IDs are the first 12 bytes of a SHA-256 over the NUL-separated
`country`, `postcode`, `city` and `street`, hex-encoded to 24 characters. See
the main [README](../README.md) for the stability guarantees.

**Note:** committed samples can lag the current code after a normalization
change. The code is the source of truth; regenerate before comparing output.
