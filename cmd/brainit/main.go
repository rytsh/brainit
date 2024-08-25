package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/rytsh/brainit"
	"github.com/rytsh/brainit/commandset"
)

var App = struct {
	CommandSet string
	File       string
}{
	CommandSet: "brainfuck",
	File:       "",
}

func argParse() {
	flag.StringVar(&App.CommandSet, "commands", App.CommandSet, "command set to interpret")
	flag.StringVar(&App.File, "file", App.File, "file to read")

	flag.Parse()
}

func main() {
	argParse()

	interpreter := brainit.NewInterpreter()

	if App.CommandSet == "brainfuck" {
		interpreter.AddCommandSet(commandset.Brainfuck)
	}

	if err := read(interpreter); err != nil {
		log.Fatal(err.Error())
	}
}

func read(interpreter *brainit.Interpreter) error {
	var reader io.Reader = os.Stdin
	if App.File != "" {
		f, err := os.Open(App.File)
		if err != nil {
			log.Fatal(err.Error())
		}

		defer f.Close()

		reader = f
	}

	return interpreter.Interpret(reader)
}
