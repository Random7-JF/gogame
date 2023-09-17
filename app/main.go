package main

import (
	"fmt"
	"os/exec"
)

func main() {
	// Define the command you want to run and its arguments
	cmd := exec.Command("ls", "-l")

	// Set up pipes for the command's standard input, output, and error
	cmd.Stdin = nil
	cmd.Stderr = nil

	// Create a pipe to capture the command's output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Read the output
	output := make([]byte, 4096)
	n, err := stdout.Read(output)
	if err != nil {
		fmt.Println("Error reading from stdout:", err)
		return
	}

	// Convert the output to a string and print it
	fmt.Println("Output:")
	fmt.Println(string(output[:n]))

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for command:", err)
		return
	}

	fmt.Println("Command finished successfully.")
}
