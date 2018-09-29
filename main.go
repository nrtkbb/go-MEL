package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/nrtkbb/go-MEL/repl"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Maya Embbeded Language!\n",
		usr.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
