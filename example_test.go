package sdl_test

import (
	"github.com/dranidis/go-sdl-spec"
	"time"
)

type UP struct {
	n int
}
type DN struct{}
type OVER struct{}

var out chan sdl.Signal

func Counter(p *sdl.Process) {
	var goingDn func() // for mutual definition

	counter := 0
	goingUp := sdl.State(p, "goingUp", func(s sdl.Signal) {
		switch v := s.(type) {
		case UP:
			counter += v.n
			if counter > 4 {
				out <- OVER{}
				defer goingDn()
				return
			}
		default:
			sdl.DefaultMessage(p, s)
		}
	})
	goingDn = sdl.State(p, "goingDn", func(s sdl.Signal) {
		switch s.(type) {
		case DN:
			counter -= 1
			if counter == 0 {
				defer goingUp()
				return
			}
		default:
			sdl.DefaultMessage(p, s)
		}
	})

	go goingUp()
}

func Example() {
	//sdl.DisableLogging()
	die := make(chan sdl.Signal)

	out = sdl.MakeBuffer()
	counterChan := sdl.MakeProcess(Counter, "Counter", die)

	go sdl.ChannelConsumer(die, "ENV", out)

	sdl.Execute(
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:UP{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:DN{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:UP{4}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:DN{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:DN{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:DN{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:DN{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:DN{}},
		sdl.Transmission{MsDelay:10, Receiver:counterChan, Signal:UP{}},
	)

	time.Sleep(2000 * time.Millisecond)
	close(die)
}
