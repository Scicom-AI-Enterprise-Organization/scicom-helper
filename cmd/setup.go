package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// setupTeleport handles the Teleport login process
func setupTeleport() error {
	fmt.Println("\n=== Teleport Setup ===")
	fmt.Println()

	// Check if already logged in
	if isTeleportLoggedIn() {
		fmt.Println("✓ You are already logged in to Teleport")

		user, err := getTeleportUser()
		if err == nil {
			fmt.Printf("✓ Logged in as: %s\n", user)
		}

		fmt.Println()
		return nil
	}

	// Not logged in, initiate login
	fmt.Println("You are not logged in to Teleport")
	fmt.Printf("Logging in to %s using GitHub SSO...\n", teleportProxy)
	fmt.Println()

	// Run tsh login interactively
	cmd := exec.Command("tsh", "login",
		fmt.Sprintf("--proxy=%s", teleportProxy),
		fmt.Sprintf("--auth=%s", githubAuth))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to log in: %v", err)
	}

	fmt.Println()
	fmt.Println("✓ Successfully logged in to Teleport!")

	// Get and display user info
	if user, err := getTeleportUser(); err == nil {
		fmt.Printf("✓ Logged in as: %s\n", user)
	}

	fmt.Println()
	return nil
}
