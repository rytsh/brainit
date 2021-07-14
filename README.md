# 🄱🅁🄰🄸🄽🄸🅃

[![Go Report Card](https://goreportcard.com/badge/github.com/rytsh/brainit?style=flat-square)](https://goreportcard.com/report/github.com/rytsh/brainit)
[![License](https://img.shields.io/github/license/rytsh/brainit?color=blue&style=flat-square)](https://raw.githubusercontent.com/rytsh/brainit/master/LICENSE)

Extensible a rune interpreter library.

Give command sets, feed interpreter with a `io.Reader` and run it.

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
    // ...

    // get new interpreter
    myInterpreter := brainit.NewInterpreter()
    // add a command set
    myInterpreter.AddCommandSet(commandset.Brainfuck)

    // give an io.Reader
    myInterpreter.Interpret(resp.Body)

    // ....
```

## License

[The MIT License (MIT)](LICENSE)

<details><summary>Tests</summary>

Test and get pprof results.

```sh
go test -cover -coverprofile cover.out -benchmem -cpuprofile cpu.out -memprofile mem.out -outputdir ./_out ./commandset/

go tool pprof commandset.test _out/cpu.out
go tool pprof commandset.test _out/mem.out

# Coverage
# Single test
# go test -cover -coverprofile cover.out -outputdir ./_out/ ./...
# Combine tests
go test -cover -coverprofile cover.out -coverpkg=github.com/rytsh/brainit,github.com/rytsh/brainit/commandset -outputdir ./_out/ ./...
go tool cover -html=./_out/cover.out
```

</details>
