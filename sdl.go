// Package sdl is an attempt to simulate SDL specifications
// in GO.
package sdl

import (
	"github.com/fatih/color"
	"sync"
	"time"
)

var mux sync.Mutex

var enterStateColor = color.New(color.FgBlue)
var nextSignalColor = color.New(color.FgRed)
var consumerColor = color.New(color.FgYellow, color.Bold)
var transmissionColor = color.New(color.FgGreen, color.Bold)

var bufferSize = 100

var logging = true

// DisableLogging turns logging off.
func DisableLogging() {
	logging = false
}

// EnableLogging turns logging on.
func EnableLogging() {
	logging = true
}

// Signal is the main structure communicated on channels.
// It can be any type.
type Signal interface{}

// Process is a type encapsulating the buffer of a process and the name of
// a process.
type Process struct {
	buffer       chan Signal
	name         string
	die          chan Signal
	saved        []Signal
	nextSaved    []Signal
	currentState string
}

func save(p *Process, s Signal) {
	p.nextSaved = append(p.nextSaved, s)
}

// DieChannel returns the channel for the termination of the process.
func DieChannel(p *Process) chan Signal {
	return p.die
}

// MakeProcess accepts a process definition and a name.
// It also receives a signal channel used for termination.
// All processes sharing the same die channel will terminate
// when close(die) is called.
// It initializes the buffer where the process is reading from.
// It returns the buffer of the process so that other
// processes can write to it.
func MakeProcess(states func(*Process), name string, die chan Signal) chan<- Signal {
	buffer := make(chan Signal, bufferSize)
	saved := []Signal{}
	nextSaved := []Signal{}
	p := Process{buffer, name, die, saved, nextSaved, ""}
	states(&p)
	return p.buffer
}

// Ignored is a helper function to print a message for ignored (consumed) messages.
// It is placed within the default section of a switch of a state.
// It prints only when Logging is enabled.
func Ignored(p *Process, s Signal) {
	if logging {
		mux.Lock()
		d := color.New(color.FgCyan)
		d.Printf("PROCESS %s AT STATE %s: IGNORES %T, %v\n", p.name, p.currentState, s, s)
		mux.Unlock()
	}
}

// State function receives the process
// and a callback function on the signal to be consumed from the buffer.
// It returns a function that will be called by the process
// when entering the state.
// State returns when the channel Done is closed.
func State(p *Process, name string, f func(s Signal)) func() {
	return func() {
		p.currentState = name
		// copy the saved signals to the actual buffer
		p.saved = make([]Signal, len(p.nextSaved))
		copy(p.saved, p.nextSaved)
		p.nextSaved = []Signal{}

		if logging {
			mux.Lock()
			enterStateColor.Printf("PROCESS %s entered STATE %s\n", p.name, p.currentState)
			mux.Unlock()
		}
		// first handle all messages in the saved buffer
		for _, s := range p.saved {
			f(s)
		}
		for { // while in this state
			s, exit := nextSignal(p)
			if exit {
				return
			}
			f(s)
		}
	}
}

func nextSignal(p *Process) (Signal, bool) {
	select {
	case s := <-p.buffer: // blocking if empty buffer
		if logging {
			mux.Lock()
			nextSignalColor.Printf("PROCESS %s AT STATE %s: %T, %v\n", p.name, p.currentState, s, s)
			mux.Unlock()
		}
		return s, false
	case <-p.die: // signal for process termination
		return nil, true
	}
}

// ChannelConsumer reads all signals at channel p and logs them at std out
// together with the name of the consumer
func ChannelConsumer(die chan Signal, n string, p chan Signal) {
	for {
		select {
		case s := <-p: // blocking if empty buffer
			mux.Lock()
			//if logging {
			consumerColor.Printf("\t\t\t\t\t%T , %v -> %s\n", s, s, n)
			mux.Unlock()
			//}
		case <-die: // signal for process termination
			return
		}
	}
}

// SendSignalsWithDelay sends all the signals in the signal list to channel c
// with a delay between each transmission equal to ms milliseconds
func SendSignalsWithDelay(c chan<- Signal, ss []Signal, ms time.Duration) {
	for _, s := range ss {
		c <- s

		mux.Lock()
		transmissionColor.Printf("%T %v\n", s, s)
		mux.Unlock()

		time.Sleep(ms * time.Millisecond)
	}
}

// MakeBuffer creates and returns a buffer for asynchronous communication
// Buffersize is defined by SetBufferSize
func MakeBuffer() chan Signal {
	return make(chan Signal, bufferSize)
}

// SetBufferSize sets the size of process buffers. Default is 100
func SetBufferSize(s int) {
	bufferSize = s
}

// Transmission is used for simulations. Defines a delay in ms, after which the signal is sent to
// the receiver channel. Executed with the Execute method for one Trasmission or with the
// Execute function for a variant number of Transmissions.
type Transmission struct {
	MsDelay  int
	Receiver chan<- Signal
	Signal   Signal
}

// Execute exetutes a number of Transmissions.
func Execute(ts ...Transmission) {
	for _, t := range ts {
		t.Execute()
	}
}

// Execute executes a single Transmission.
func (t Transmission) Execute() {
	time.Sleep(time.Duration(t.MsDelay) * time.Millisecond)
	t.Receiver <- t.Signal

	mux.Lock()
	transmissionColor.Printf("%T %v\n", t.Signal, t.Signal)
	mux.Unlock()
}

// DefaultMessage is a helper function for printing a message that it is consumed as
// a default action at a switch signal.
func DefaultMessage(p *Process, s Signal) {
	mux.Lock()
	d := color.New(color.FgCyan)
	d.Printf("------ At state %s ignored %v, %T\n", p.currentState, s, s)
	mux.Unlock()
}
