package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var countWorkers = 20
var countJobs = 200
var countErrors = 40

func startWorker(workerNum int, in <-chan func() error, chError chan error, quit chan bool, wg *sync.WaitGroup) {
	for input := range in {
		select {
		case <-quit:
			wg.Done()
		default:
			chError <- input()
		}
	}
}
func startError(in chan func() error, chError chan error, quit chan bool, wg *sync.WaitGroup) {
	var count int
	for err := range chError {
		fmt.Printf("Eror text %v \n", err)
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
	runtime.GOMAXPROCS(4)
	var wg sync.WaitGroup

	workerInput := make(chan func() error, 0)
	chError := make(chan error, 0)
	quit := make(chan bool)

	var sliceWork []func() error
	for i := 1; i <= countJobs; i++ {
		work := work(i)
		sliceWork = append(sliceWork, work)
	}

	go startError(workerInput, chError, quit, &wg)

	for z := 0; z < countWorkers; z++ {
		go startWorker(z, workerInput, chError, quit, &wg)
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
