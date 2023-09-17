package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

type MineServerInstance struct {
	Path   string
	Jar    string
	Name   string
	Input  chan string
	Output chan string
}

func main() {
	go startServer("/home/random/minecraft/server1/", "server.jar", "Server1")
	go startServer("/home/random/minecraft/server2/", "server.jar", "Server2")
}

func startServer(serverPath string, serverJar string, instanceName string) error {
	javaPath, err := exec.LookPath("java")
	if err != nil {
		log.Println(instanceName + ": Java not found in Path")
		return err
	}
	// Start the Minecraft server
	cmd := exec.Command(javaPath, "-Xmx1024M", "-Xms1024M", "-jar", (serverPath + serverJar), "nogui")
	cmd.Dir = serverPath
	// Create pipes for stdin and stdout
	stdinPipe, _ := cmd.StdinPipe()
	stdoutPipe, _ := cmd.StdoutPipe()

	// Start the server
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}

	// Create channels for communication between goroutines
	fromStdin := make(chan string)
	fromStdout := make(chan string)

	// Goroutine for reading from stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			fromStdin <- text
		}
	}()

	// Goroutine for writing to stdin
	go func() {
		for {
			select {
			case input := <-fromStdin:
				io.WriteString(stdinPipe, input+"\n")
			}
		}
	}()

	// Goroutine for reading from stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			text := scanner.Text()
			fromStdout <- text
		}
	}()

	// Goroutine for writing to stdout (or do any processing on the output)
	go func() {
		for {
			select {
			case output := <-fromStdout:
				fmt.Println(instanceName+"-Output:", output)
			}
		}
	}()

	// Wait for the server to finish (you can remove this if you want to keep the server running)
	err = cmd.Wait()
	if err != nil {
		fmt.Println(instanceName+"Error waiting for server:", err)
		return err
	}
	return nil
}
