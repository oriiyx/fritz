package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
)

// CalculateSchemaHash computes SHA256 hash of an EntityDefinition struct
// Uses canonical JSON serialization to ensure consistent hashing
func CalculateSchemaHash(definition *definitions.EntityDefinition) (string, error) {
	// Convert struct to canonical JSON bytes
	jsonBytes, err := json.Marshal(definition)
	if err != nil {
		return "", fmt.Errorf("failed to marshal definition to JSON: %w", err)
	}

	// Calculate SHA256 hash
	hash := sha256.Sum256(jsonBytes)

	// Convert to hex string
	return hex.EncodeToString(hash[:]), nil
}

// CalculateJSONHash computes SHA256 hash directly from JSON bytes
// Useful when you already have the JSON representation
func CalculateJSONHash(jsonBytes []byte) string {
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:])
}

// CalculateFileHash computes SHA256 hash from a JSON file
// Reads the file and calculates hash in one operation
func CalculateFileHash(filePath string) (string, error) {
	jsonBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return CalculateJSONHash(jsonBytes), nil
}

// CompareSchemaHashes is a helper to quickly check if two definitions differ
func CompareSchemaHashes(def1, def2 *definitions.EntityDefinition) (bool, error) {
	hash1, err := CalculateSchemaHash(def1)
	if err != nil {
		return false, fmt.Errorf("failed to hash first definition: %w", err)
	}

	hash2, err := CalculateSchemaHash(def2)
	if err != nil {
		return false, fmt.Errorf("failed to hash second definition: %w", err)
	}

	return hash1 == hash2, nil
}
