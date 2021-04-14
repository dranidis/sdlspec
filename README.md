# SDL specification and simulation in GO

[![GoDoc](https://godoc.org/github.com/dranidis/sdlspec?status.svg)](https://godoc.org/github.com/dranidis/sdlspec)

The purpose of this package is to allow easy definition and simulation of SDL (Specification and Description Language http://sdl-forum.org/index.htm) processes in GO.

This is an experimental project and is not fully tested yet.

## Hello world example

```go
package main

import (
	"fmt"
	"time"

	"github.com/dranidis/sdlspec"
)

type HI struct{}

func helloStates(p *sdlspec.Process) {
	start := sdlspec.State(p, "start", func(s sdlspec.Signal) {
		switch s.(type) {
		case HI:
			fmt.Println("Hello SDL")
		default:
		}
	})
	go start()
}

func main() {
	die := make(chan sdlspec.Signal)
	helloProcess := sdlspec.MakeProcess(helloStates, "hello", die)
	helloProcess <- HI{}

	time.Sleep(1000 * time.Millisecond)
	close(die)
}
```

The output is:

```
PROCESS hello AT STATE start: main.HI, {}
Hello SDL

```

## Signals

Each SDL signal is declared as a new go struct type:

```go
type HI struct {}
```

## Processes

A process is created using the sdlspec.MakeProcess function:

```go
	helloProcess := sdlspec.MakeProcess(helloStates, "hello", die)
```

that takes as a parameter a function like the following:

```go
func helloStates(p *sdlspec.Process) {
	start := sdlspec.State(p, "start", func(s sdlspec.Signal) {
		switch s.(type) {
		case HI:
			fmt.Println("Hello SDL")
		default:
		}
	})
	go start()
}
```

The function defines a state **start** using the construction:

```go
	start := sdlspec.State(p, "start", func(s sdlspec.Signal) { ... })
```

The callback function defines the behaviour at that state. The important part is within the switch statement:

```go
		switch s.(type) {
		case HI:
			fmt.Println("Hello SDL")
		default:
		}
```

If the received signal is of type HI, print the Hello SDL message. Else ignore the signal. Note that the signal is consumed anyway.

The process is spawned at the initial state with:

```go
	go start()
```

In the main:

```go
    helloProcess <- HI{}
    time.Sleep(2000 * time.Millisecond)
    close(die)
```

we send the signal `HI{}` to the process, sleep for 2 secs and terminate all SDL processes by closing the `die` channel.
