package main

import (
	"log"
	"net/http"

	"github.com/rytsh/brainit"
	"github.com/rytsh/brainit/commandset"
)

func main() {
	resp, err := http.Get("https://raw.githubusercontent.com/erikdubbelboer/brainfuck-jit/master/mandelbrot.bf")
	if err != nil {
		log.Fatalln("Upps:", err)
	}
	defer resp.Body.Close()

	myInterpreter := brainit.NewInterpreter()
	myInterpreter.AddCommandSet(commandset.Brainfuck)

	myInterpreter.Interpret(resp.Body)
}
