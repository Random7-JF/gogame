package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"

	"github.com/Random7-JF/gogame/app/server"
)

func main() {

	javaBin, err := exec.LookPath("java")
	if err != nil {
		log.Println("Java not found in Path")
		panic(err)
	}

	server1 := server.Init("Server1", "/home/random/minecraft/server1/", "server.jar", javaBin)
	server2 := server.Init("Server2", "/home/random/minecraft/server2/", "server.jar", javaBin)
	currentInstance := server1

	server1.Start()
	server2.Start()

	var wg sync.WaitGroup
	wg.Add(1)

	for {
		test := <-server1.FromStdOut
		fmt.Println(test)
	}

	scanner := bufio.NewScanner(os.Stdin)
	go func() {
		for scanner.Scan() {
			text := scanner.Text()
			fmt.Println("Select server instance or Q to quit")
			switch text {
			case "1":
				if currentInstance != server1 {
					currentInstance = server1
				}
			case "2":
				if currentInstance != server2 {
					currentInstance = server2
				}
			case "q":
				wg.Done()
				return
			default:
				currentInstance.FromStdIn <- text
			}
		}
	}()

	wg.Wait()
}
