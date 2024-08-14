package main

import (
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

/*
	UnBuffered channel won't let the writer write until
	one reader is ready to read from the channel,

	therefore the following solution only work, where prcessB
	waits on the channel before sending to release channel.


	-------------------------------------------------implementation I----------------------------------------------------

	func processA(wait <-chan struct{}, release chan<- struct{}) {
	fmt.Println("in process A")
	release <- struct{}{}
	fmt.Println("waiting for B's signal")
	<-wait
	fmt.Println("Got Signal from B")
	fmt.Println("process A is executing")
	}

	func processB(wait <-chan struct{}, release chan<- struct{}) {
		fmt.Println("in process B")
		fmt.Println("waiting for A's signal")
		<-wait
		release <- struct{}{}
		fmt.Println("Got Signal from A")
		fmt.Println("process B is executing")
	}

	func main() {
		ch1 := make(chan struct{})
		ch2 := make(chan struct{})

		go processA(ch1, ch2)
		go processB(ch2, ch1)

		time.Sleep(time.Second * 10)
	}


	-------------------------------------------------------------------------------------------------------------------
	if process B follows other order, i.e.,
	release()
	wait()
	it will cause a deadlock

	--------------------------------------------------deadlock---------------------------------------------------------

	func processA(wait <-chan struct{}, release chan<- struct{}) {
	fmt.Println("in process A")
	release <- struct{}{}
	fmt.Println("waiting for B's signal")
	<-wait
	fmt.Println("Got Signal from B")
	fmt.Println("process A is executing")
	}

	func processB(wait <-chan struct{}, release chan<- struct{}) {
		fmt.Println("in process B")
		release <- struct{}{}
		fmt.Println("waiting for A's signal")
		<-wait
		fmt.Println("Got Signal from A")
		fmt.Println("process B is executing")
	}

	func main() {
		ch1 := make(chan struct{})
		ch2 := make(chan struct{})

		go processA(ch1, ch2)
		go processB(ch2, ch1)

		time.Sleep(time.Second * 10)
	}

	-------------------------------------------------------------------------------------------------------------------

	one way to avoid this is to use buffered channel, implemented below
*/

/*-------------------------------------------------Implementation II---------------------------------------------------*/

func processA(wait <-chan struct{}, release chan<- struct{}) {
	fmt.Println("process A is executing")
	release <- struct{}{}
	fmt.Println("waiting for B's signal")
	<-wait
	fmt.Println("Got Signal from B")
	time.Sleep(5 * time.Second)
	fmt.Println("execution DONE : process A")
}

func processB(wait <-chan struct{}, release chan<- struct{}) {
	fmt.Println("process B is executing")
	release <- struct{}{}
	fmt.Println("waiting for A's signal")
	<-wait
	fmt.Println("Got Signal from A")
	time.Sleep(5 * time.Second)
	fmt.Println("execution DONE : process B")
}

func main() {
	ch1 := make(chan struct{}, 1) // can be unbuffered for this impl
	ch2 := make(chan struct{}, 1)

	errgrp := errgroup.Group{}

	errgrp.Go(func() error {
		processB(ch2, ch1)
		return nil
	})
	errgrp.Go(func() error {
		processA(ch1, ch2)
		return nil
	})

	if err := errgrp.Wait(); err != nil {
		fmt.Printf("error while processing : %s", err.Error())
	}
}
