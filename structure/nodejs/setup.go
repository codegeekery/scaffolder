package main

import (
	"fmt"
	"os"
	"os/exec"
)


func main() {
	fmt.Println("Setting up SvelteKit dependencies...")

	cmd := exec.Command("npm", "install")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: Failed to install SvelteKit dependencies.")
		os.Exit(1)
	}

	_ = os.Remove("setup.go")

	fmt.Println("Setup completed successfully!")
}
