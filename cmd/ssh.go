package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/AlecAivazis/survey/v2"
)

// sshToNode allows user to select a node and SSH into it
func sshToNode() error {
	fmt.Println("\n=== Teleport SSH ===")
	fmt.Println()

	// Check if logged in
	if !isTeleportLoggedIn() {
		fmt.Println("You are not logged in to Teleport")
		fmt.Println("Please run 'Teleport Setup' first")
		return fmt.Errorf("not logged in to Teleport")
	}

	// Get list of nodes
	fmt.Println("Fetching available nodes...")
	nodes, err := getTeleportNodes()
	if err != nil {
		return fmt.Errorf("failed to get nodes: %v", err)
	}

	if len(nodes) == 0 {
		fmt.Println("No nodes available")
		return nil
	}

	fmt.Printf("Found %d node(s)\n\n", len(nodes))

	// Let user select a node
	var selectedNode string
	prompt := &survey.Select{
		Message: "Select a node to connect to:",
		Options: nodes,
		PageSize: 15,
	}

	err = survey.AskOne(prompt, &selectedNode)
	if err != nil {
		return fmt.Errorf("selection cancelled")
	}

	// Get available logins for the selected node
	fmt.Printf("\nTesting available logins for %s...\n", selectedNode)
	fmt.Println("(This may take a few seconds)")
	logins, err := getNodeLogins(selectedNode)
	if err != nil {
		return fmt.Errorf("failed to get logins: %v", err)
	}

	if len(logins) == 0 {
		return fmt.Errorf("no logins available for this node")
	}

	// Let user select a login
	var selectedLogin string
	loginPrompt := &survey.Select{
		Message: "Select a login user:",
		Options: logins,
		Default: logins[0], // Default to first login (usually ubuntu)
	}

	err = survey.AskOne(loginPrompt, &selectedLogin)
	if err != nil {
		return fmt.Errorf("selection cancelled")
	}

	// Connect to the selected node with the selected login
	fmt.Printf("\nConnecting to %s as %s...\n", selectedNode, selectedLogin)
	fmt.Println("(Press Ctrl+D or type 'exit' to disconnect)")
	fmt.Println()

	// Run tsh ssh interactively with login user
	cmd := exec.Command("tsh", "ssh", fmt.Sprintf("%s@%s", selectedLogin, selectedNode))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Don't treat normal exit as an error
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 130 { // Ctrl+C
				fmt.Println("\nConnection closed")
				return nil
			}
		}
		return fmt.Errorf("SSH connection failed: %v", err)
	}

	fmt.Println("\nConnection closed")
	return nil
}
