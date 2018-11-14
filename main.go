package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/nrtkbb/go-MEL/lexer"
	"github.com/nrtkbb/go-MEL/parser"
	"github.com/nrtkbb/go-MEL/repl"
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		// REPL mode
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Maya Embbeded Language!\n",
			usr.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	filePaths := flag.Args()
	for _, fp := range filePaths {
		if _, err := os.Stat(fp); err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		fmt.Println(fp)
		input, err := ioutil.ReadFile(fp)
		if err != nil {
			log.Fatal(err)
		}

		l := lexer.New(string(input))
		p := parser.New(l)
		program := p.ParseProgram()
		fmt.Println(program.String())
	}
}
