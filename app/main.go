package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/Random7-JF/gogame/app/server"
)

func main() {
	var current *server.MinecraftServer

	javaBin, err := exec.LookPath("java")
	if err != nil {
		log.Println("Java not found in Path")
		panic(err)
	}

	server1 := server.Init("Server1", "/home/random/minecraft/server1/", "server.jar", javaBin)
	server1.Start()

	server2 := server.Init("Server2", "/home/random/minecraft/server2/", "server.jar", javaBin)
	server2.Start()

	current = server1

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			switch text {
			case "1":
				if current != server1 {
					current = server1
				}
			case "2":
				if current != server2 {
					current = server2
				}
			case "q":
				current.Stop()
				wg.Done()
				return
			default:
				current.FromStdIn <- text
			}
		}
	}()

	wg.Wait()
}
