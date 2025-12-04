package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// isTshInstalled checks if the tsh command is available
func isTshInstalled() bool {
	cmdName := "which"
	if runtime.GOOS == "windows" {
		cmdName = "where"
	}
	cmd := exec.Command(cmdName, "tsh")
	return cmd.Run() == nil
}

// isTeleportLoggedIn checks if the user is logged in to Teleport
func isTeleportLoggedIn() bool {
	cmd := exec.Command("tsh", "status")
	return cmd.Run() == nil
}

// getTeleportUser returns the current Teleport user
func getTeleportUser() (string, error) {
	cmd := exec.Command("tsh", "status")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Parse output to find User: line
	// Format can be "User: username" or "> Profile URL:  user@cluster"
	output := out.String()
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Try format: "User: username"
		if strings.HasPrefix(trimmed, "User:") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}

		// Try format: "> Profile URL:  https://proxy/web/cluster/user@cluster" or similar
		if strings.Contains(trimmed, "Profile URL:") && strings.Contains(trimmed, "@") {
			// Extract user from URL format
			atIndex := strings.LastIndex(trimmed, "@")
			if atIndex > 0 {
				// Look backwards from @ to find the username
				beforeAt := trimmed[:atIndex]
				lastSlash := strings.LastIndex(beforeAt, "/")
				if lastSlash >= 0 {
					user := beforeAt[lastSlash+1:]
					if user != "" {
						return user, nil
					}
				}
			}
		}

		// Try format: "Logged in as: username"
		if strings.Contains(trimmed, "Logged in as:") {
			parts := strings.Split(trimmed, ":")
			if len(parts) >= 2 {
				user := strings.TrimSpace(parts[1])
				if user != "" {
					return user, nil
				}
			}
		}
	}

	return "", fmt.Errorf("could not find user in tsh status output. Output was:\n%s", output)
}

// getTeleportNodes returns a list of available Teleport nodes
func getTeleportNodes() ([]string, error) {
	cmd := exec.Command("tsh", "ls", "--format=names")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	nodes := []string{}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			nodes = append(nodes, line)
		}
	}

	return nodes, nil
}

// pickDefaultLogin picks the best default login from available logins
// Priority: ubuntu > root > first available
func pickDefaultLogin(logins []string) string {
	if len(logins) == 0 {
		return "ubuntu"
	}

	// Check for ubuntu first
	for _, login := range logins {
		if login == "ubuntu" {
			return "ubuntu"
		}
	}

	// Check for root second
	for _, login := range logins {
		if login == "root" {
			return "root"
		}
	}

	// Return first available
	return logins[0]
}

// getAllLogins returns all logins from tsh status (all roles combined)
func getAllLogins() ([]string, error) {
	cmd := exec.Command("tsh", "status")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	output := out.String()
	lines := strings.Split(output, "\n")

	// Look for "Logins:" line
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "Logins:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) >= 2 {
				loginsStr := strings.TrimSpace(parts[1])
				loginList := strings.Split(loginsStr, ",")
				logins := []string{}
				for _, login := range loginList {
					login = strings.TrimSpace(login)
					if login != "" {
						logins = append(logins, login)
					}
				}
				return logins, nil
			}
		}
	}

	return nil, fmt.Errorf("no logins found")
}

// testLoginAccess tests if a specific login works for a node
func testLoginAccess(login, nodeName string) bool {
	// Use a quick SSH test command with short timeout
	// The "exit" command will close immediately if login is allowed
	cmd := exec.Command("tsh", "ssh", fmt.Sprintf("%s@%s", login, nodeName), "exit")

	// Set a timeout context (2 seconds should be enough)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()

	// If command succeeds, login is allowed
	// If it fails with "access denied" or similar, login is not allowed
	if err != nil {
		// Check if error is due to access denial
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "access denied") ||
		   strings.Contains(stderrStr, "permission denied") ||
		   strings.Contains(stderrStr, "not allowed") {
			return false
		}
		// Other errors might be temporary, so we'll include the login
		return true
	}

	return true
}

// getNodeLogins returns logins that actually work for the specific node
func getNodeLogins(nodeName string) ([]string, error) {
	// Get all possible logins from tsh status
	allLogins, err := getAllLogins()
	if err != nil {
		// Fall back to defaults if we can't get logins
		return []string{"ubuntu", "root"}, nil
	}

	// Test each login to see which ones work for this node
	validLogins := []string{}

	for _, login := range allLogins {
		if testLoginAccess(login, nodeName) {
			validLogins = append(validLogins, login)
		}
	}

	// If no logins passed the test, return all of them anyway
	// (the test might have failed for other reasons)
	if len(validLogins) == 0 {
		return allLogins, nil
	}

	return validLogins, nil
}

// runCommand executes a command and returns its output
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	return out.String(), nil
}

// runInteractiveCommand runs a command with stdin/stdout/stderr attached
func runInteractiveCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}
