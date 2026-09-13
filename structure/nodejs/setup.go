package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("Setting up project dependencies...")

	cmd := exec.Command("npm", "install", "express", "cors")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		fmt.Println("Error: Failed to install dependencies.")
		os.Exit(1)
	}

	// Se borra a sí mismo para no dejar rastro en la carpeta final del usuario
	_ = os.Remove("setup.go")

	fmt.Println("Setup completed successfully!")
}
