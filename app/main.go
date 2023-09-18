package main

import (
	"log"
	"os/exec"

	"github.com/Random7-JF/gogame/app/server"
)

func main() {
	var err error
	var server1 server.MinecraftInstance
	server1.Name = "Server1"
	server1.Jar = "server.jar"
	server1.Path = "/home/random/minecraft/server1/"

	var server2 server.MinecraftInstance
	server2.Name = "Server2"
	server2.Jar = "server.jar"
	server2.Path = "/home/random/minecraft/server2/"

	server1.JavaBin, err = exec.LookPath("java")
	if err != nil {
		log.Println("Java not found in Path")
		panic(err)
	}

	server2.JavaBin, err = exec.LookPath("java")
	if err != nil {
		log.Println("Java not found in Path")
		panic(err)
	}

	go server1.StartServer()
	go server2.StartServer()
	select {}

}
