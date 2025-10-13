//
// Copyright (c) 2025 Open2b Software Snc. All Rights Reserved.
//
// WARNING: This software is protected by international copyright laws.
// Redistribution in part or in whole strictly prohibited.
//

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var connectorsInfoJSONFilename = filepath.Join("..", "catalog.json")

func main() {

	f, err := os.Open(connectorsInfoJSONFilename)
	if err != nil {
		fatal("cannot open JSON file: %s", err)
	}
	defer f.Close()

	var jsonFile map[string]any
	err = json.NewDecoder(f).Decode(&jsonFile)
	if err != nil {
		fatal("cannot decode JSON file: %s", err)
	}

	// Validate the JSON file.
	err = validateJSONFile(jsonFile)
	if err != nil {
		fatal("validation of JSON file failed: %s", err)
	}
	fmt.Printf("✅ Declaration of connectors is valid\n")

	// Format the JSON file.
	err = formatJSON(jsonFile, connectorsInfoJSONFilename)
	if err != nil {
		fatal("error while formatting the JSON file: %s", err)
	}
	fmt.Printf("✅ JSON file reformatted\n")

}

func fatal(format string, a ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}
