package main

import (
	"fmt"
	"sync"
	"time"
)

var countWorkers = 20
var countJobs = 40
var countErrors = 20

func startWorker(in <-chan func() error, chError chan error, quit chan bool, wg *sync.WaitGroup) {
	for input := range in {
		select {
		case <-quit:
			wg.Done()
		default:
			chError <- input()
		}
	}
}
func startError(chError chan error, quit chan bool, wg *sync.WaitGroup) {
	var count int
	for err := range chError {
		if err != nil {
			fmt.Printf("Eror text %v \n", err)
		}
		count++
		if count == countErrors {
			fmt.Println("after stop")
			close(quit)
		}
		wg.Done()
	}
}

func main() {
	if countWorkers > countJobs {
		countWorkers = countJobs
	}
	if countErrors > countJobs {
		countErrors = countJobs
	}
	var wg sync.WaitGroup

	workerInput := make(chan func() error)
	chError := make(chan error)
	quit := make(chan bool)

	var sliceWork []func() error
	for i := 1; i <= countJobs; i++ {
		work := work(i)
		sliceWork = append(sliceWork, work)
	}

	go startError(chError, quit, &wg)

	for z := 0; z < countWorkers; z++ {
		go startWorker(workerInput, chError, quit, &wg)
	}

	for _, f := range sliceWork {
		wg.Add(1)
		select {
		case <-quit:
			wg.Done()
		default:
			workerInput <- f
		}
	}
	close(workerInput)
	wg.Wait()
}

func work(x int) func() error {
	return func() error {
		for i := 1; i > 0; i++ {
			if i%x*1000 == 0 {
				time.Sleep(1000 * time.Millisecond)
				return fmt.Errorf("fail work: %d", x)
			}
		}
		return nil
	}
}
