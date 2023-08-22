package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

func main() {
	stackFolder := "stack" // Folder containing the subdirectories with YAML files

	actionPrompt := promptui.Select{
		Label: "Select an action:",
		Items: []string{"Create stack", "Remove stack"},
	}

	_, action, err := actionPrompt.Run()
	if err != nil {
		fmt.Println("Action selection failed:", err)
		return
	}

	switch action {
	case "Create stack":
		createStack(stackFolder)
	case "Remove stack":
		removeStack(stackFolder)
	default:
		fmt.Println("Invalid action selected.")
	}
}

func createStack(stackFolder string) {
	stackDirs := listSubdirectories(stackFolder)

	if len(stackDirs) == 0 {
		fmt.Println("No subdirectories found in the folder.")
		return
	}

	var options []string
	for _, stackDir := range stackDirs {
		baseDir := filepath.Base(stackDir)
		options = append(options, baseDir)
	}

	prompt := promptui.Select{
		Label: "Select a stack to create:",
		Items: options,
	}

	selectedIdx, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	selectedDir := stackDirs[selectedIdx]
	selectedFile := filepath.Join(selectedDir, "docker-compose.yaml")

	// Execute docker-compose up -d for the selected YAML file
	cmd := exec.Command("docker-compose", "-f", selectedFile, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error starting containers:", err)
		return
	}

	fmt.Println("Stack started successfully!")
}

func removeStack(stackFolder string) {
	stackDirs := listSubdirectories(stackFolder)

	if len(stackDirs) == 0 {
		fmt.Println("No subdirectories found in the folder.")
		return
	}

	var options []string
	for _, stackDir := range stackDirs {
		baseDir := filepath.Base(stackDir)
		options = append(options, baseDir)
	}

	prompt := promptui.Select{
		Label: "Select a stack to remove:",
		Items: options,
	}

	selectedIdx, _, err := prompt.Run()
	if err != nil {
		fmt.Println("Prompt failed:", err)
		return
	}

	selectedDir := stackDirs[selectedIdx]

	// Execute docker-compose down for the selected stack
	cmd := exec.Command("docker-compose", "-f", filepath.Join(selectedDir, "docker-compose.yaml"), "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		fmt.Println("Error removing stack:", err)
		return
	}

	fmt.Println("Stack removed successfully!")
}

func listSubdirectories(folderPath string) []string {
	var subDirs []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() != folderPath {
			if _, err := os.Stat(filepath.Join(path, "docker-compose.yaml")); err == nil {
				subDirs = append(subDirs, path)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error listing subdirectories:", err)
	}

	return subDirs
}
