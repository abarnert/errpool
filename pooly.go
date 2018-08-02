package main

// Different ways of limiting concurrency:
// 1. Build a semaphore out of a channel with max length.
//    (Can abstract with a semaphore type, but it really doesn't
//    add much--sem := semaphore.New(32) is a bit nicer than
//    sem := make(chan bool, 32), etc., but you still have to
//    make sure you get the qcquire, release, and pre-load right.)
//  1.1. Wait on the semaphore at the start of each goroutine.
//       (Creates an unbounded number of goroutines, but this doesn't
//       seem to be an actual performance problem as you'd expect.)
//  1.2. Wait on the semaphore before launching each goroutine
//       (Requires starting the semaphore at full instead of empty,
//       but that doesn't make things that much more complicated.)
// 2. Build a pool of max length goroutines that service a channel.
//    (Can easily abstract with a wrapper around errgroup, unless
//    I'm missing something. Could also extend it to, e.g., allow
//    separate shutdown and wait methods, but I didn't need any
//    such extensions here, so I stuck with the errgroup API.)

import (
	"context"
	"./errpool"
	"time"

	"golang.org/x/sync/errgroup"
)

func doit() (int, error) {
	for i := 0; i < 1000; i++ {
		time.Sleep(2 * time.Microsecond)
	}
	return 1, nil
}

func doitSpam(n int) error {
	g, _ := errgroup.WithContext(context.Background())
	results := make(chan int, n)
	for i := 0; i < n; i++ {
		g.Go(func() error {
			res, err := doit()
			if err != nil {
				return err
			}
			results <- res
			return nil
		})
	}
	return g.Wait()
}

func doitSemaphore(n int) error {
	g, _ := errgroup.WithContext(context.Background())
	sem := make(chan bool, 32)
	defer close(sem)
	results := make(chan int, n)
	for i := 0; i < n; i++ {
		sem <- true
		g.Go(func() error {
			defer func() { <-sem }()
			res, err := doit()
			if err != nil {
				return err
			}
			results <- res
			return nil
		})
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	for i := 0; i < n; i++ {
		<-results
	}
	return g.Wait()
}

func doitSemaphoreBlock(n int) error {
	g, _ := errgroup.WithContext(context.Background())
	sem := make(chan bool, 32)
	defer close(sem)
	results := make(chan int, n)
	for i := 0; i < n; i++ {
		g.Go(func() error {
			sem <- true
			defer func() { <-sem }()
			res, err := doit()
			if err != nil {
				return err
			}
			results <- res
			return nil
		})
	}
	return g.Wait()
}

func doitPool(n int) error {
	g, _ := errgroup.WithContext(context.Background())
	jobs := make(chan func() (int, error), n)
	results := make(chan int, n)
	for i := 0; i < n; i++ {
		jobs <- doit
	}
	close(jobs)
	for i := 0; i < 32; i++ {
		g.Go(func() error {
			for f := range jobs {
				res, err := f()
				if err != nil {
					return err
				}
				results <- res
			}
			return nil
		})
	}
	for i := 0; i < n; i++ {
		<-results
	}
	return g.Wait()
}

func doitErrPool(n int) error {
	p, _ := errpool.WithContext(context.Background(), 32)
	results := make(chan int, n)
	for i := 0; i < n; i++ {
		p.Go(func() error {
			res, err := doit()
			if err != nil {
				return err
			}
			results <- res
			return nil
		})
	}
	for i := 0; i < n; i++ {
		<-results
	}
	return p.Wait()
}

func main() {
	doitSemaphore(10)
	doitPool(10)
}
