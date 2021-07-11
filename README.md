# brainit

Extensible a rune interpreter library.

Feed interpreter with a `io.Reader` and run it.  
Give own functions and add.

__Command Sets__  
__-__ Brainfuck

## Usage

```go
import (
    "github.com/rytsh/brainit"
    "github.com/rytsh/brainit/commandset"
)

    // ...
    resp, err := http.Get("https://raw.githubusercontent.com/erikdubbelboer/brainfuck-jit/master/mandelbrot.bf")
    if err != nil {
        log.Fatalln("Upps:", err)
    }

    // get new interpreter
    myInterpreter := brainit.NewInterpreter()
    // add a command set
    myInterpreter.AddCommandSet(commandset.Brainfuck)

    // give an io.Reader
    myInterpreter.Interpret(resp.Body)
```

## License

[The MIT License (MIT)](LICENSE)

<details><summary>Testing</summary>

Test and get pprof results explanations.

```sh
go test -cover -coverprofile cover.out -benchmem -cpuprofile cpu.out -memprofile mem.out -outputdir ./_out ./commandset/

go tool pprof commandset.test _out/cpu.out
go tool pprof commandset.test _out/mem.out
```

</details>
