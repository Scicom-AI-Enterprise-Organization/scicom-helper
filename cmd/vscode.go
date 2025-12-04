package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// editorConfig represents an editor's settings configuration
type editorConfig struct {
	name         string
	settingsPath string
}

// getEditorSettingsPaths returns settings paths for VS Code and Cursor
func getEditorSettingsPaths(home string) []editorConfig {
	var editors []editorConfig

	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			editors = []editorConfig{
				{name: "VS Code", settingsPath: filepath.Join(appData, "Code", "User", "settings.json")},
				{name: "Cursor", settingsPath: filepath.Join(appData, "Cursor", "User", "settings.json")},
			}
		}
	} else if runtime.GOOS == "darwin" {
		editors = []editorConfig{
			{name: "VS Code", settingsPath: filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")},
			{name: "Cursor", settingsPath: filepath.Join(home, "Library", "Application Support", "Cursor", "User", "settings.json")},
		}
	} else {
		// Linux
		editors = []editorConfig{
			{name: "VS Code", settingsPath: filepath.Join(home, ".config", "Code", "User", "settings.json")},
			{name: "Cursor", settingsPath: filepath.Join(home, ".config", "Cursor", "User", "settings.json")},
		}
	}

	return editors
}

// configureEditor configures a specific editor for Teleport compatibility
func configureEditor(editor editorConfig) error {
	// Check if editor settings directory exists
	editorDir := filepath.Dir(editor.settingsPath)
	if _, err := os.Stat(editorDir); os.IsNotExist(err) {
		return fmt.Errorf("%s not found (directory doesn't exist)", editor.name)
	}

	// Read existing settings or create empty map
	var settings map[string]interface{}
	data, err := os.ReadFile(editor.settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create new settings
			settings = make(map[string]interface{})
			fmt.Printf("Creating new %s settings file...\n", editor.name)
		} else {
			return fmt.Errorf("failed to read %s settings: %v", editor.name, err)
		}
	} else {
		// Parse existing settings
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to parse %s settings: %v", editor.name, err)
		}
		fmt.Printf("Found existing %s settings\n", editor.name)
	}

	// Check if setting already exists and is correct
	if val, exists := settings["remote.SSH.useLocalServer"]; exists {
		if boolVal, ok := val.(bool); ok && !boolVal {
			fmt.Printf("✓ %s already configured correctly for Teleport\n", editor.name)
			fmt.Println("  remote.SSH.useLocalServer = false")
			return nil
		}
	}

	// Backup existing settings
	if len(data) > 0 {
		backupPath := editor.settingsPath + ".backup"
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

	if err := os.WriteFile(editor.settingsPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write %s settings: %v", editor.name, err)
	}

	fmt.Printf("✓ %s configured successfully\n", editor.name)
	fmt.Println("  Set remote.SSH.useLocalServer = false")

	return nil
}

// configureVSCode sets up VS Code and Cursor settings for Teleport compatibility
func configureVSCode() error {
	fmt.Println("\n=== Configure Editors for Teleport ===")
	fmt.Println()

	// Get home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Get editor settings paths
	editors := getEditorSettingsPaths(home)

	// Track which editors we successfully configured
	configuredEditors := []string{}
	notFoundEditors := []string{}

	// Try to configure each editor
	for _, editor := range editors {
		fmt.Printf("Checking for %s...\n", editor.name)
		if err := configureEditor(editor); err != nil {
			if os.IsNotExist(err) || filepath.Dir(editor.settingsPath) == "" {
				notFoundEditors = append(notFoundEditors, editor.name)
				fmt.Printf("  %s not found (not installed or not run yet)\n", editor.name)
			} else {
				fmt.Printf("  Warning: Failed to configure %s: %v\n", editor.name, err)
			}
		} else {
			configuredEditors = append(configuredEditors, editor.name)
		}
		fmt.Println()
	}

	// Print summary
	if len(configuredEditors) == 0 {
		return fmt.Errorf("no editors found. Please install VS Code or Cursor and run it at least once")
	}

	fmt.Println("=== Editor Configuration Complete! ===")
	fmt.Println()
	fmt.Printf("✓ Configured %d editor(s):\n", len(configuredEditors))
	for _, name := range configuredEditors {
		fmt.Printf("  - %s\n", name)
	}
	fmt.Println()
	fmt.Println("This setting is required for Teleport SSH connections.")
	fmt.Println("Restart your editor(s) for changes to take effect.")
	fmt.Println()

	return nil
}
