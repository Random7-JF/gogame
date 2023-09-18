package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
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
	server.Instance = exec.Command(javaBin, "-Xmx1024M", "-Xms1024M", "-jar", (path + jar), "nogui")
	server.Instance.Dir = path
	server.StdIn, _ = server.Instance.StdinPipe()
	server.StdOut, _ = server.Instance.StdoutPipe()

	return &server
}

func (m *MinecraftServer) Start() {
	err := m.Instance.Start()
	if err != nil {
		log.Println("Instance Start - Error starting server: ", err)
		return
	}
	go m.ReadFromStdOut()
	go m.WriteToStdIn()
}

func (m *MinecraftServer) Stop() {
	err := m.StdIn.Close()
	if err != nil {
		log.Println("StdIn Close - Error closing: ", err)
		return
	}
	err = m.Instance.Wait()
	if err != nil {
		log.Println("Instance Wait - Error closing: ", err)
		return
	}
}

// for writing to stdin CLI -> Minecraft
func (m *MinecraftServer) WriteToStdIn() {
	for {
		select {
		//input is value of what in channel.
		case input := <-m.FromStdIn:
			//write input with new line.
			io.WriteString(m.StdIn, input+"\n")
		}
	}
}

// for reading from stdout Minecraft -> CLI
func (m *MinecraftServer) ReadFromStdOut() {
	scanner := bufio.NewScanner(m.StdOut)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Server %s output: %s\n", m.Name, text)
		m.FromStdOut <- text
	}
}
