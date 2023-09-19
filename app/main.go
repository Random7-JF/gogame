package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/Random7-JF/gogame/app/server"
)

func main() {
	var current *server.MinecraftServer

	javaBin, err := exec.LookPath("java")
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(javaBin, "-Xmx1024M", "-Xms1024M", "-jar", serverPath, "nogui")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = serverName

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	//	scanner := bufio.NewScanner(os.Stdin)

	for {
		test := <-server1.FromStdOut
		fmt.Println(test)
	}

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		defer stdin.Close()
		for {
			select {
			case text := <-input:
				fmt.Fprintln(stdin, text)
			}
		}
	}()

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	server1Input := make(chan string)
	server2Input := make(chan string)

	go startServer("/home/random/minecraft/server1/server.jar", "/home/random/minecraft/server1", server1Input)
	go startServer("/home/random/minecraft/server2/server.jar", "/home/random/minecraft/server2", server2Input)

	for {
		var serverName string
		var command string

		fmt.Print("Enter server name (server1/server2) and command: ")
		fmt.Scan(&serverName, &command)

		switch serverName {
		case "server1":
			server1Input <- command
		case "server2":
			server2Input <- command
		default:
			fmt.Println("Invalid server name")
		}
	}
}
