package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

var (
	AccessSecret  string
	RefreshSecret string
	ConfigSecret  string
)

func generateRandomHex(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func InitEnv() error {
	envPath := filepath.Join("./db", ".env")

	// If .env doesn't exist, create it with random secrets
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		access := generateRandomHex(32)
		refresh := generateRandomHex(32)
		configSec := generateRandomHex(32)

		content := "access_secret=" + access + "\n" +
			"refresh_secret=" + refresh + "\n" +
			"config_secret=" + configSec + "\n"

		if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
			return err
		}
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			switch key {
			case "access_secret":
				AccessSecret = val
			case "refresh_secret":
				RefreshSecret = val
			case "config_secret":
				ConfigSecret = val
			}
		}
	}

	// Fallback if any is empty
	if AccessSecret == "" {
		AccessSecret = generateRandomHex(32)
	}
	if RefreshSecret == "" {
		RefreshSecret = generateRandomHex(32)
	}
	if ConfigSecret == "" {
		ConfigSecret = generateRandomHex(32)
	}

	return nil
}
