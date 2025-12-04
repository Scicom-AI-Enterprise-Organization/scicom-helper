package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

const (
	teleportProxy = "teleport-iam.aies.scicom.dev"
	githubAuth    = "github-connector"
)

var rootCmd = &cobra.Command{
	Use:   "scicom-helper",
	Short: "Scicom Platform Engineering helper tool",
	Long:  `A CLI tool to help Scicom developers manage infrastructure access, starting with Teleport EC2 access.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMode()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func runInteractiveMode() {
	// Check if tsh is installed first
	if !isTshInstalled() {
		fmt.Println("Error: tsh (Teleport CLI) is not installed")
		fmt.Println("Please install from: https://goteleport.com/download")
		os.Exit(1)
	}

	for {
		var choice string
		prompt := &survey.Select{
			Message: "What would you like to do?",
			Options: []string{
				"Teleport Setup (Login)",
				"Teleport Update Nodes (Update SSH config)",
				"Teleport SSH (Connect to a node)",
				"Exit",
			},
		}

		err := survey.AskOne(prompt, &choice)
		if err != nil {
			fmt.Println("Exiting...")
			return
		}

		switch choice {
		case "Teleport Setup (Login)":
			if err := setupTeleport(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Teleport Update Nodes (Update SSH config)":
			if err := updateNodes(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Teleport SSH (Connect to a node)":
			if err := sshToNode(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Exit":
			fmt.Println("Goodbye!")
			return
		}

		fmt.Println()
	}
}
