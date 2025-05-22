package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type state int

var GlobalTimer int
var IsPaused bool
var GlobalState state

var WORKING_TIME int
var RESTING_TIME int

const (
	RESTING state = iota
	WORKING
	NOTHING
)

type Pomodoro struct {
	max_time     int
	time_elapsed chan int
}

func (p *Pomodoro) UpdateSeconds() {
	for {
		if !IsPaused || GlobalState == NOTHING {
			// need to run infinite times until globalstate enum changes to smth else
			if GlobalTimer > p.max_time && GlobalState != NOTHING {
				p.ToggleState()
			}
			p.time_elapsed <- GlobalTimer
			time.Sleep(time.Millisecond * 1000)
			GlobalTimer++
		}
	}
}

func (p *Pomodoro) ChangeState(maxTime int, status state) {
	p.max_time = maxTime
	GlobalState = status
	// reset time_elapsed when changing state
	GlobalTimer = 1
}

func (p *Pomodoro) ToggleState() {
	if GlobalState == WORKING {
		p.ChangeState(RESTING_TIME, RESTING)
		return
	}
	p.ChangeState(WORKING_TIME, WORKING)
	return
}

func main() {
	working_time := flag.Int("w", 60*120, "long working time in second")
	resting_time := flag.Int("r", 60*30, "short resting time in second")
	flag.Parse()
	WORKING_TIME = *working_time
	RESTING_TIME = *resting_time

	GlobalTimer = 1
	IsPaused = false
	elapsed := make(chan int, 1)
	GlobalState = NOTHING

	toggle_state := make(chan os.Signal, 1)
	toggle_pause := make(chan os.Signal, 1)

	signal.Notify(toggle_state, syscall.SIGUSR1)
	signal.Notify(toggle_pause, syscall.SIGUSR2)

	pomo := Pomodoro{
		max_time:     WORKING_TIME,
		time_elapsed: elapsed,
	}

	go pomo.UpdateSeconds()

	go func() {
		for {
			<-toggle_state
			// we don't want to wait 1 second waiting for elapsed to be sent from UpdateSeconds
			// waybar will not show our program for 1 second which looks werid and switching state
			// will incur 1second delay otherwise
			elapsed <- 1
			pomo.ToggleState()
		}
	}()
	go func() {
		for {
			<-toggle_pause
			IsPaused = !IsPaused
		}
	}()

	for val := range elapsed {
		switch GlobalState {
		case RESTING:
			fmt.Printf("{\"text\": \"%v\", \"alt\": \"RESTING\"}\n", time.Duration(val*1e9))
		case WORKING:
			fmt.Printf("{\"text\": \"%v\", \"alt\": \"WORKING\"}\n", time.Duration(val*1e9))
		case NOTHING:
			fmt.Printf("{\"text\": \"Pomo\", \"alt\": \"NON\"}\n")
		}
	}
}
