package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

// setupTeleportGitHub handles the Teleport login process using GitHub SSO
func setupTeleportGitHub() error {
	fmt.Println("\n=== Teleport Setup (GitHub SSO) ===")
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

	// Run tsh login interactively with GitHub auth
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

// setupTeleportLocal handles the Teleport login process using local account
func setupTeleportLocal() error {
	fmt.Println("\n=== Teleport Setup (Local Account) ===")
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
	fmt.Printf("Logging in to %s using local account...\n", teleportProxy)
	fmt.Println()

	// Prompt for username
	fmt.Print("Enter your Teleport username: ")
	var username string
	fmt.Scanln(&username)

	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	fmt.Println("You will be prompted for your password.")
	fmt.Println()

	// Run tsh login interactively with local auth and specified username
	cmd := exec.Command("tsh", "login",
		fmt.Sprintf("--proxy=%s", teleportProxy),
		fmt.Sprintf("--user=%s", username),
		"--auth=local")
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
