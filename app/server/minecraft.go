package server

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type MinecraftServer struct {
	Path       string
	Jar        string
	Name       string
	JavaBin    string
	Instance   *exec.Cmd
	StdIn      io.WriteCloser
	StdOut     io.ReadCloser
	FromStdIn  chan string
	FromStdOut chan string
}

func Init(name string, path string, jar string, javaBin string) *MinecraftServer {
	var server MinecraftServer
	server.Name = name
	server.Path = path
	server.Jar = jar
	server.JavaBin = javaBin
	// Start the Minecraft server
	server.Instance = exec.Command(javaBin, "-Xmx1024M", "-Xms1024M", "-jar", (path + jar), "nogui")
	server.Instance.Dir = path
	// Create pipes for stdin and stdout
	server.StdIn, _ = server.Instance.StdinPipe()
	server.StdOut, _ = server.Instance.StdoutPipe()

	return &server
}

func (m *MinecraftServer) Start() {
	// Start the server
	err := m.Instance.Start()
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	go m.readFromStdOut()
	go m.writeToStdIn()

}

func (m *MinecraftServer) Stop() {
	m.StdIn.Close()
	m.Instance.Wait()
}

func (m *MinecraftServer) StartIO() {
	go m.writeToStdIn()
	go m.readFromStdOut()
}

func (m *MinecraftServer) StopIO() {
	m.StdIn.Close()

}

// Goroutine for writing to stdin
func (m *MinecraftServer) writeToStdIn() {
	for {
		select {
		case input := <-m.FromStdIn:
			io.WriteString(m.StdIn, input+"\n")
		}
	}
}

// Goroutine for reading from stdout
func (m *MinecraftServer) readFromStdOut() {
	scanner := bufio.NewScanner(m.StdOut)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Server %d output: %s\n", m.Name, text)

		m.FromStdOut <- text
	}
}
