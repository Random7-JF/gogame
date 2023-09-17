package main

import (
	"os"
	"os/exec"
)

func main() {
	// Define the command you want to run and its arguments
	err := serve()
	if err != nil {
		panic(err)
	}

}

func serve() error {
	log, err := os.Create("output.log")
	if err != nil {
		return err
	}
	defer log.Close()
	cmd := exec.Command("bash", " tmux new -s vaulthunters ./serve.sh")
	cmd.Stdout = log
	cmd.Stderr = log
	return cmd.Run()
}
