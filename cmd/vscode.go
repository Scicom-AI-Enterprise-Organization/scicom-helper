package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// configureVSCode sets up VS Code settings for Teleport compatibility
func configureVSCode() error {
	fmt.Println("\n=== Configure VS Code for Teleport ===")
	fmt.Println()

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// VS Code settings path (cross-platform)
	var vscodeSettingsPath string
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return fmt.Errorf("APPDATA environment variable not set")
		}
		vscodeSettingsPath = filepath.Join(appData, "Code", "User", "settings.json")
	} else if runtime.GOOS == "darwin" {
		vscodeSettingsPath = filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
	} else {
		// Linux
		vscodeSettingsPath = filepath.Join(home, ".config", "Code", "User", "settings.json")
	}

	// Check if VS Code settings directory exists
	vscodeDir := filepath.Dir(vscodeSettingsPath)
	if _, err := os.Stat(vscodeDir); os.IsNotExist(err) {
		fmt.Println("VS Code settings directory not found.")
		fmt.Println("Make sure VS Code is installed and has been run at least once.")
		return fmt.Errorf("VS Code not found")
	}

	// Read existing settings or create empty map
	var settings map[string]interface{}
	data, err := os.ReadFile(vscodeSettingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create new settings
			settings = make(map[string]interface{})
			fmt.Println("Creating new VS Code settings file...")
		} else {
			return fmt.Errorf("failed to read VS Code settings: %v", err)
		}
	} else {
		// Parse existing settings
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to parse VS Code settings: %v", err)
		}
		fmt.Println("Found existing VS Code settings")
	}

	// Check if setting already exists and is correct
	if val, exists := settings["remote.SSH.useLocalServer"]; exists {
		if boolVal, ok := val.(bool); ok && !boolVal {
			fmt.Println("✓ VS Code already configured correctly for Teleport")
			fmt.Println("  remote.SSH.useLocalServer = false")
			return nil
		}
	}

	// Backup existing settings
	if len(data) > 0 {
		backupPath := vscodeSettingsPath + ".backup"
		if err := os.WriteFile(backupPath, data, 0644); err != nil {
			fmt.Printf("Warning: Failed to create backup: %v\n", err)
		} else {
			fmt.Printf("Backed up existing settings to: %s\n", backupPath)
		}
	}

	// Update setting
	settings["remote.SSH.useLocalServer"] = false

	// Write updated settings
	updatedData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %v", err)
	}

	if err := os.WriteFile(vscodeSettingsPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write VS Code settings: %v", err)
	}

	fmt.Println()
	fmt.Println("=== VS Code Configuration Complete! ===")
	fmt.Println()
	fmt.Println("✓ Set remote.SSH.useLocalServer = false")
	fmt.Println()
	fmt.Println("This setting is required for Teleport SSH connections.")
	fmt.Println("Restart VS Code for changes to take effect.")
	fmt.Println()

	return nil
}
