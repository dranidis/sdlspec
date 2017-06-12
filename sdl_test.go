package sdl

import (
	"testing"
	"time"
)


type HI struct{}
type HO struct{}


var testVar int = 0

func helloStates(p *Process) {
	start := State(p, func(s Signal) {
		switch s.(type) {
		case HI:
			testVar = 1
		default:
		}
	})
	go start()
}

func TestSendRecieve(t *testing.T) {
	helloProcess := MakeProcess(helloStates, "hello")

	helloProcess <- HI{}
	time.Sleep(100 * time.Millisecond)

	if testVar != 1 {
		t.Error("case HI not executed")
	}
}

func TestEndProcesses(t *testing.T) {
	SetBufferSize(50)
	
	helloProcess := MakeProcess(helloStates, "hello2")

	SendSignalsWithDelay(helloProcess, []Signal{
		HI{}, HI{}, HI{},
	}, 10)	

	// Done should be encapsulated into the process
	//EndProcesses()
	time.Sleep(500 * time.Millisecond)
	if testVar != 1 {
		t.Error("case HI not executed 2")
	}

}

func TestChannelConsumer(t *testing.T) {
	out := MakeBuffer()

	helloStates1 := func (p *Process) {
		start := State(p, func(s Signal) {
			switch s.(type) {
			case HI:
				out <- HO{}
				testVar = 1
			default:
			}
		})
		go start()
	}

	helloProcess := MakeProcess(helloStates1, "hello3")
	go ChannelConsumer("OUT", out)
	time.Sleep(1000 * time.Millisecond)

	SendSignalsWithDelay(helloProcess, []Signal{
		HI{}, HI{}, HI{},
	}, 10)	

	time.Sleep(500 * time.Millisecond)


	EndProcesses()
	time.Sleep(500 * time.Millisecond)
	if testVar != 1 {
		t.Error("case HI not executed 3")
	}

}