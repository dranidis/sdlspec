# go-sdl-spec

The purpose of this package is to allow easy definition and simulation of SDL (Specification and Description Language http://sdl-forum.org/index.htm) processes in GO.

This is an experimental project and is not fully tested yet.

## Hello world example

```go
package main

import (
	"fmt"
	"github.com/dranidis/go-sdl-spec"
	"time"
)

type HI struct{}

func helloStates(buffer sdl.SignalChan) {
	start := sdl.State(buffer, func(s sdl.Signal) {
		switch s.(type) {
		case HI:
			fmt.Println("Hello SDL")
		default:
		}
	})
		go start()
}

func main() {
	helloProcess := sdl.Process(helloStates)

	helloProcess <- HI{}

	time.Sleep(2000 * time.Millisecond)
	sdl.EndProcesses()
}
```

The output is:
```
PROCESS:  main.HI, {}
Hello SDL
```

## Signals
Each SDL signal is declared as a new go struct type:
```go
type HI struct {}
```
## Processes
A process is created using the sdl.Process function:
```go
	helloProcess := sdl.Process(helloStates)
```
that takes as a parameter a function like the following:
```go
func helloStates(buffer sdl.SignalChan) {
	start := sdl.State(buffer, func(s sdl.Signal) {
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
	start := sdl.State(buffer, func(s sdl.Signal) { ... })
```
The callback function defines the behaviour at that state. The important part is within the switch statement:
```go
		case HI:
			fmt.Println("Hello SDL")
		default:
```
If the received signal is of type HI, print the Hello SDL message. Else ignore it.

The process is spawned at the initial state with:
```go
		go start()
```
In the main:
```go
    helloProcess <- HI{}
    time.Sleep(2000 * time.Millisecond)
    sdl.EndProcesses()
```
we send the signal `HI{}` to the process, sleep for 2 secs and terminate all SDL processes.