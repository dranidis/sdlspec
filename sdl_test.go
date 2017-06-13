package sdl

import (
	"testing"
	"time"
)

type HI struct{}
type HO struct{}

func TestSendRecieve(t *testing.T) {
	testVar := 0
	die := make(chan Signal)
	helloProcess := MakeProcess(func(p *Process) {
		start := State(p, func(s Signal) {
			switch s.(type) {
			case HI:
				testVar = 1
			default:
			}
		})
		go start()
	}, "hello", die)

	helloProcess <- HI{}
	time.Sleep(100 * time.Millisecond)

	if testVar != 1 {
		t.Error("case HI not executed")
	}
	// EndProcesses()
	close(die)

}

func TestEndProcesses(t *testing.T) {
	testVar := 0
	die := make(chan Signal)

	SetBufferSize(50)

	helloProcess := MakeProcess(func(p *Process) {
		start := State(p, func(s Signal) {
			switch s.(type) {
			case HI:
				testVar = 1
			default:
			}
		})
		go start()
	}, "hello1", die)

	SendSignalsWithDelay(helloProcess, []Signal{
		HI{}, HI{}, HI{},
	}, 10)

	// Done should be encapsulated into the process
	//EndProcesses()
	time.Sleep(500 * time.Millisecond)
	if testVar != 1 {
		t.Error("case HI not executed 2")
	}
	close(die)
}

func TestChannelConsumer(t *testing.T) {
	testVar := 0
	die := make(chan Signal)

	out := MakeBuffer()

	helloStates1 := func(p *Process) {
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

	helloProcess := MakeProcess(helloStates1, "hello3", die)
	go ChannelConsumer(die, "OUT", out)
	time.Sleep(1000 * time.Millisecond)

	SendSignalsWithDelay(helloProcess, []Signal{
		HI{}, HI{}, HI{},
	}, 10)

	time.Sleep(500 * time.Millisecond)

	//EndProcesses()
	time.Sleep(500 * time.Millisecond)
	if testVar != 1 {
		t.Error("case HI not executed 3")
	}
	close(die)

}
