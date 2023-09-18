package server

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type MinecraftInstance struct {
	Path    string
	Jar     string
	Name    string
	JavaBin string
	StdIn   io.WriteCloser
	StdOut  io.ReadCloser
	Input   chan string
	Output  chan string
}

func (m *MinecraftInstance) StartServer() error {
	// Start the Minecraft server
	cmd := exec.Command(m.JavaBin, "-Xmx1024M", "-Xms1024M", "-jar", (m.Path + m.Jar), "nogui")
	cmd.Dir = m.Path
	// Create pipes for stdin and stdout
	m.StdIn, _ = cmd.StdinPipe()
	m.StdOut, _ = cmd.StdoutPipe()

	// Start the server
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return err
	}

	// Create channels for communication between goroutines
	m.Input = make(chan string)
	m.Output = make(chan string)

	go m.readFromStdIn()
	go m.writeToStdIn()
	go m.readFromStdOut()
	go m.writeToStdOut()

	// Wait for the server to finish (you can remove this if you want to keep the server running)
	err = cmd.Wait()
	if err != nil {
		fmt.Println(m.Name+"Error waiting for server:", err)
		return err
	}
	return nil
}

// Goroutine for writing to stdin
func (m *MinecraftInstance) writeToStdIn() {
	for {
		select {
		case input := <-m.Input:
			io.WriteString(m.StdIn, input+"\n")
		}
	}
}

// Goroutine for reading from stdin
func (m *MinecraftInstance) readFromStdIn() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		m.Input <- text
	}
}

// Goroutine for reading from stdout
func (m *MinecraftInstance) readFromStdOut() {
	scanner := bufio.NewScanner(m.StdOut)
	for scanner.Scan() {
		text := scanner.Text()
		m.Output <- text
	}
}

// Goroutine for writing to stdout (or do any processing on the output)
func (m *MinecraftInstance) writeToStdOut() {
	for {
		select {
		case output := <-m.Output:
			fmt.Println(m.Name+"-Output:", output)
		}
	}
}
