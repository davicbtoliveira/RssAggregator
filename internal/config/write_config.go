package config

import (
	"encoding/json"
	"os"
)

func write(cfg *Config) error {
	byteArray, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	os.WriteFile(filePath, byteArray, 0644)
	return nil
}
