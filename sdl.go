// An attempt to simulate SDL specifications
// in GO.
package sdl

import (
	"fmt"
	"time"
)

var bufferSize = 100

// The channel Done is used for terminating all processes.
// Each process returns after a select on channel Done.
// The main closes Done at termination by calling EndProcess.
var Done chan bool = make(chan bool)

// Signal is the main structure communicated on channels.
// It can be any type.
type Signal interface{}
type SignalChan chan Signal

// Process function accepts a process definition.
// It initializes the buffer where the process is reading from.
// It returns the buffer of the process so that other
// processes can write to it.
func Process(states func(SignalChan)) chan Signal {
	buffer := make(chan Signal, bufferSize)
	states(buffer)
	return buffer
}

// State function receives the buffer of the process
// and a callback function on the signal to be consumed from the buffer.
// It returns a function that will be called by the process
// when entering the state.
// State returns when the channel Done is closed.
func State(buffer SignalChan, f func(s Signal)) func() {
	return func() {
		for { // while in this state
			s, exit := nextSignal(buffer)
			if exit {
				return
			}
			f(s)
		}
	}
}

func nextSignal(b SignalChan) (Signal, bool) {
	select {
	case s := <-b: // blocking if empty buffer
		fmt.Printf("PROCESS:  %T, %v\n", s, s)
		return s, false
	case <-Done: // signal for process termination
		return nil, true
	}
}

// Closes the Done channel so that all processes terminate.
func EndProcesses() {
	close(Done)
}

// Reads all signals at channel p and logs them at std out
// together with the name of the consumer
func ChannelConsumer(n string, p SignalChan) {
	for {
		select {
		case s := <-p: // blocking if empty buffer
			fmt.Printf("%s <- %T, %v\n", n, s, s)
		case <-Done: // signal for process termination
			return
		}
	}
}

// Sends all the signals in the signal list to channel c
// with a delay between each transmission equal to ms milliseconds
func SendSignalsWithDelay(c SignalChan, ss []Signal, ms time.Duration) {
	for _, s := range ss {
		c <- s
		time.Sleep(ms * time.Millisecond)
	}
}

// Creates and returns a buffer for asynchronous communication
// Buffersize is defined by SetBufferSize
func MakeBuffer() chan Signal {
	return make(chan Signal, bufferSize)
}

// Sets the size of process buffers. Default is 100
func SetBufferSize(s int) {
	bufferSize = s
}