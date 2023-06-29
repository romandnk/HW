package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		log.Fatalf("too many arguments")
	}

	path := args[1]

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf(`directory by path "%s" is not exist`, path)
		}
		log.Fatal(err)
	}

	env, err := ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(args[2:], env)
	os.Exit(returnCode)
}
