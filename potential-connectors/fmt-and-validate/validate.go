//
// Copyright (c) 2025 Open2b Software Snc. All Rights Reserved.
//
// WARNING: This software is protected by international copyright laws.
// Redistribution in part or in whole strictly prohibited.
//

package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func validateJSONFile(jsonFile map[string]any) error {

	// Check that the map does not contain any unexpected keys.
	if unexpected := unexpectedMapKeys(jsonFile, []string{"connectors"}); unexpected != nil {
		return fmt.Errorf("JSON object has unexpected keys: %v", strings.Join(unexpected, ", "))
	}

	connectors, ok := jsonFile["connectors"].([]any)
	if !ok {
		return errors.New("the JSON file must contain a key 'connectors' of type JSON array")
	}

	for _, connector := range connectors {
		connector, ok := connector.(map[string]any)
		if !ok {
			return errors.New("every connector must be a JSON object")
		}
		err := validateConnector(connector)
		if err != nil {
			name, ok := connector["name"].(string)
			if ok && name != "" {
				return fmt.Errorf("invalid connector %q: %s", name, err)
			}
			return err
		}
	}

	return nil

}

func validateConnector(connector map[string]any) error {

	// Check that the map does not contain any unexpected keys.
	allowedMapKeys := []string{"code", "label", "categories", "connectorType", "asSource", "asDestination"}
	if unexpected := unexpectedMapKeys(connector, allowedMapKeys); unexpected != nil {
		return fmt.Errorf("JSON object representing connector has unexpected keys: %v", strings.Join(unexpected, ", "))
	}

	// Field 'code'.
	code, ok := connector["code"].(string)
	if !ok || code == "" {
		return errors.New("field 'code' must be a not-empty JSON string")
	}

	// Field 'label'.
	label, ok := connector["label"].(string)
	if !ok || label == "" {
		return errors.New("field 'label' must be a not-empty JSON string")
	}

	// Field 'categories'.
	categories, ok := connector["categories"].([]any)
	if !ok || len(categories) == 0 {
		return errors.New("field 'categories' must be a not-empty JSON array of not-empty strings")
	}
	for _, category := range categories {
		category, ok := category.(string)
		if !ok || category == "" {
			return errors.New("field 'categories' must be a not-empty JSON array of not-empty strings")
		}
	}

	// Field 'connectorType'.
	connectorType, ok := connector["connectorType"].(string)
	if !ok {
		return errors.New("missing field 'connectorType'")
	}
	switch connectorType {
	case "API", "Database", "File", "FileStorage", "MessageBroker", "SDK", "Webhook":
		// Ok.
	default:
		return fmt.Errorf("invalid value for 'connectorType': %q", connectorType)
	}

	var hasAsSource, hasAsDestination bool

	// Field 'asSource', if present.
	if asSource, ok := connector["asSource"]; ok {
		asSource, ok := asSource.(map[string]any)
		if !ok {
			return errors.New("field 'asSource' must be a JSON Object")
		}
		err := validateAsRole(asSource)
		if err != nil {
			return fmt.Errorf("invalid JSON Object for 'asSource': %s", err)
		}
		hasAsSource = true
	}

	// Field 'asDestination', if present.
	if asSource, ok := connector["asDestination"]; ok {
		asSource, ok := asSource.(map[string]any)
		if !ok {
			return errors.New("field 'asDestination' must be a JSON Object")
		}
		err := validateAsRole(asSource)
		if err != nil {
			return fmt.Errorf("invalid JSON Object for 'asDestination': %s", err)
		}
		hasAsDestination = true
	}

	// At least one key between 'asSource' and 'asDestination' must be
	// present.
	if !hasAsSource && !hasAsDestination {
		return errors.New("at least one key between 'asSource' and 'asDestination' must be present")
	}

	return nil

}

func validateAsRole(asRole map[string]any) error {

	// Check that the map does not contain any unexpected keys.
	if unexpected := unexpectedMapKeys(asRole, []string{"description", "implemented", "comingSoon"}); unexpected != nil {
		return fmt.Errorf("JSON object has unexpected keys: %v", strings.Join(unexpected, ", "))
	}

	description, ok := asRole["description"].(string)
	if !ok || description == "" {
		return errors.New("field 'description' must be provided and must be a not-empty JSON string")
	}

	// Field 'implemented'.
	implemented, ok := asRole["implemented"].(bool)
	if !ok {
		return errors.New("field 'implemented' must be provided and must be a JSON boolean")
	}

	// Field 'comingSoon'.
	comingSoon, ok := asRole["comingSoon"].(bool)
	if !ok {
		return errors.New("field 'comingSoon' must be provided and must be a JSON boolean")
	}

	if implemented && comingSoon {
		return errors.New("cannot have both 'implemented' and 'comingSoon' set to true")
	}

	return nil

}

// unexpectedMapKeys validates the keys of the map m, returning the keys that
// are not among the allowed ones.
func unexpectedMapKeys[T any](m map[string]T, allowedKeys []string) (unexpected []string) {
	for k := range m {
		if !slices.Contains(allowedKeys, k) {
			unexpected = append(unexpected, k)
		}
	}
	return unexpected
}
