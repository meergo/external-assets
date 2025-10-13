//
// Copyright (c) 2025 Open2b Software Snc. All Rights Reserved.
//
// WARNING: This software is protected by international copyright laws.
// Redistribution in part or in whole strictly prohibited.
//

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

func formatJSON(jsonFile map[string]any, toFile string) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetIndent("", "    ")
	enc.SetEscapeHTML(false)

	type asRoleType struct {
		Description string `json:"description"`
		Implemented bool   `json:"implemented"`
		ComingSoon  bool   `json:"comingSoon"`
	}

	type cType struct {
		Code          string      `json:"code"`
		Label         string      `json:"label"`
		Categories    []any       `json:"categories"`
		ConnectorType string      `json:"connectorType"`
		AsSource      *asRoleType `json:"asSource,omitempty"`
		AsDestination *asRoleType `json:"asDestination,omitempty"`
	}
	var toSerial []cType
	for _, connector := range jsonFile["connectors"].([]any) {
		connector := connector.(map[string]any)
		c := cType{}
		c.Code = connector["code"].(string)
		c.Label = connector["label"].(string)
		c.Categories = connector["categories"].([]any)
		c.ConnectorType = connector["connectorType"].(string)
		if asSource, ok := connector["asSource"].(map[string]any); ok {
			c.AsSource = &asRoleType{
				Description: asSource["description"].(string),
				Implemented: asSource["implemented"].(bool),
				ComingSoon:  asSource["comingSoon"].(bool),
			}
		}
		if asDestination, ok := connector["asDestination"].(map[string]any); ok {
			c.AsDestination = &asRoleType{
				Description: asDestination["description"].(string),
				Implemented: asDestination["implemented"].(bool),
				ComingSoon:  asDestination["comingSoon"].(bool),
			}
		}
		toSerial = append(toSerial, c)
	}
	err := enc.Encode(map[string]any{"connectors": toSerial})
	if err != nil {
		return fmt.Errorf("cannot encode data into JSON: %s", err)
	}
	err = os.WriteFile(toFile, b.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("cannot write file %q: %s", toFile, err)
	}
	return nil
}
