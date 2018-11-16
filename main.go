package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/nrtkbb/go-MEL/lexer"
	"github.com/nrtkbb/go-MEL/parser"
	"github.com/nrtkbb/go-MEL/repl"
)

var mel = regexp.MustCompile(`.mel$`)

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
		stat, err := os.Stat(fp)
		if err != nil && !os.IsExist(err) {
			log.Println(err)
			continue
		}

		if stat.IsDir() {
			err = readDir(fp)
			log.Println(err)
			continue
		}

		if !mel.MatchString(fp) {
			continue
		}

		readFile(fp)
	}
}

func readDir(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !mel.MatchString(path) {
			return nil
		}

		err = readFile(path)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func readFile(file string) error {
	fmt.Println(file)
	input, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	l := lexer.New(string(input))
	p := parser.New(l)
	program := p.ParseProgram()
	fmt.Println(program.String())

	return nil
}
