package osm2addr

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// writeJsonFile ...
func writeJsonFile(countrycode, filename string, in map[string]*TagSET) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}

// writeJsonFileMap ...
func writeJsonFileMap(countrycode, filename string, in map[string]bool) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}

// writeJsonFileMMap ...
func writeJsonFileMMap(countrycode, filename string, in map[string]map[string]bool) {
	folder := filepath.Join("json", countrycode)
	_ = os.MkdirAll(folder, 0755)
	j, err := json.MarshalIndent(&in, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
		panic(err)
	}
}

// writeJsonFileIDMap ...
//func writeJsonFileIDMap(countrycode, filename string, in map[string]objectID) {
//	folder := filepath.Join("json", countrycode)
//	_ = os.MkdirAll(folder, 0755)
//	out := make(map[string]string, len(in))
//	for key, value := range in {
//		out[key] = value.b64()
//	}
//	j, err := json.MarshalIndent(&out, "", "\t")
//	if err != nil {
//		panic(err)
//	}
//	if err := os.WriteFile(filepath.Join(folder, filename), j, 0644); err != nil {
//		panic(err)
//	}
//}
