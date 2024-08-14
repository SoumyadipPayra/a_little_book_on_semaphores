package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func processA(ctx context.Context, waitSemaphore, releaseSemaphore *semaphore.Weighted) {
	fmt.Println("Executing Process A")
	releaseSemaphore.Release(0)
	fmt.Println("waiting for signal from Process B")
	waitSemaphore.Acquire(ctx, 0)
	fmt.Println("Got sigal from process B")
	time.Sleep(time.Second * 5)
	fmt.Println("Done execution : process A")
}

func processB(ctx context.Context, waitSemaphore, releaseSemaphore *semaphore.Weighted) {
	fmt.Println("Executing Process B")
	releaseSemaphore.Release(0)
	fmt.Println("waiting for signal from Process A")
	waitSemaphore.Acquire(ctx, 0)
	fmt.Println("Got sigal from process A")
	time.Sleep(time.Second * 5)
	fmt.Println("Done execution : process B")
}

func main() {
	sem1 := semaphore.NewWeighted(0)
	sem2 := semaphore.NewWeighted(0)

	ctx := context.TODO()

	errgrp := errgroup.Group{}
	errgrp.Go(func() error { processA(ctx, sem1, sem2); return nil })
	errgrp.Go(func() error { processB(ctx, sem2, sem1); return nil })

	if err := errgrp.Wait(); err != nil {
		fmt.Printf("error while procesing : %s", err.Error())
	}

}
