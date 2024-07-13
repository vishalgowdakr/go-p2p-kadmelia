package cmd

import (
	"encoding/base64"
	"os"
)

func EncodeNodeID(nodeID string) string {
	return base64.StdEncoding.EncodeToString([]byte(nodeID))
}

func WriteToFile(encodedNodeID, filePath string) error {
	return os.WriteFile(filePath, []byte(encodedNodeID), 0444) // Read-only permissions
}

func ReadFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func DecodeNodeID(encodedNodeID string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedNodeID)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
