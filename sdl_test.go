package sdlspec

import (
	"testing"
	"time"
)

type HI struct {
	n int
}
type HO struct {
	n int
}

func TestCaseAndDefault(t *testing.T) {
	testVar := 0
	die := make(chan Signal)
	helloProcess := MakeProcess(func(p *Process) {
		start := State(p, "start", func(s Signal) {
			switch s.(type) {
			case HI:
				testVar = 1
			default:
				testVar = 2
			}
		})
		go start()
	}, "hello", die)

	helloProcess <- HI{}
	time.Sleep(100 * time.Millisecond)

	if testVar != 1 {
		t.Error("case HI not executed")
	}
	helloProcess <- HO{}
	time.Sleep(100 * time.Millisecond)

	if testVar != 2 {
		t.Error("case default not executed")
	}
	// EndProcesses()
	close(die)

}

func TestChangingState(t *testing.T) {
	testVar := 0
	die := make(chan Signal)

	SetBufferSize(50)

	helloProcess := MakeProcess(func(p *Process) {
		var next func()
		start := State(p, "start", func(s Signal) {
			switch s.(type) {
			case HI:
				defer next()
				return
			default:
			}
		})
		next = State(p, "next", func(s Signal) {
			switch s.(type) {
			case HI:
				testVar++
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
	if testVar != 2 {
		t.Error("case HI at state next not executed")
	}
	close(die)
}

func TestSaveSignal(t *testing.T) {
	testVar := 0
	die := make(chan Signal)

	SetBufferSize(50)

	helloProcess := MakeProcess(func(p *Process) {
		var next func()
		start := State(p, "start", func(s Signal) {
			switch s.(type) {
			case HI:
				defer next()
				return
			case HO:
				save(p, s)
			default:
			}
		})
		next = State(p, "next", func(s Signal) {
			switch s.(type) {
			case HO:
				testVar++
			default:
			}
		})
		go start()
	}, "hello1", die)

	SendSignalsWithDelay(helloProcess, []Signal{
		HO{1}, HO{2}, HO{3}, HI{4}, HI{5},
	}, 10)

	// Done should be encapsulated into the process
	//EndProcesses()
	time.Sleep(500 * time.Millisecond)
	if testVar != 3 {
		t.Error("signal HO not saved at state start")
	}
	close(die)
}

func TestSaveThenIgnoreSignal(t *testing.T) {
	testVar := 0
	die := make(chan Signal)

	SetBufferSize(50)

	helloProcess := MakeProcess(func(p *Process) {
		var next1, next2 func()
		start := State(p, "start", func(s Signal) {
			switch s.(type) {
			case HI:
				defer next1()
				return
			case HO:
				save(p, s)
			default:
			}
		})
		next1 = State(p, "next1", func(s Signal) {
			switch s.(type) {
			case HI:
				defer next2()
				return
			default:
			}
		})
		next2 = State(p, "next2", func(s Signal) {
			switch s.(type) {
			case HO:
				testVar++
			case HI:
			default:
			}
		})
		go start()
	}, "hello2", die)

	SendSignalsWithDelay(helloProcess, []Signal{
		HO{1}, HO{2}, HO{3}, HI{4}, HI{5}, HI{6},
	}, 10)

	// Done should be encapsulated into the process
	//EndProcesses()
	time.Sleep(500 * time.Millisecond)
	if testVar != 0 {
		t.Error("signal HO not forgotten at state next1")
	}
	close(die)
}

func TestChannelConsumer(t *testing.T) {
	testVar := 0
	die := make(chan Signal)

	out := MakeBuffer()

	helloStates1 := func(p *Process) {
		start := State(p, "start", func(s Signal) {
			switch s.(type) {
			case HI:
				out <- HO{}
				testVar = -3
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
	if testVar != -3 {
		t.Error("case HI not executed 3")
	}
	close(die)

}
