package main

import (
	"fmt"
	"runtime"
	"sync"
)

var countWorkers = 200
var countJobs = 2000
var countErrors = 400

func startWorker(workerNum int, in <-chan func() error, chError chan error, quit chan bool, wg *sync.WaitGroup) {
	fmt.Printf("start worker %d \n", workerNum)
	for input := range in {
		select {
		case <-quit:
			wg.Done()
			fmt.Println("stop worker")
		default:
			fmt.Println("start job")
			chError <- input()
		}
		//wg.Done()
		fmt.Println("stop job")
	}
	fmt.Println("worker all job done")
}
func startError(in chan func() error, chError chan error, quit chan bool, wg *sync.WaitGroup) {
	var count int
	for err := range chError {
		fmt.Printf("Eror text %v \n", err)
		count++
		if count == countErrors {
			fmt.Println("close channels")
			close(quit)
			//wg.Done()
			//fmt.Println("- break")
			//return
		}
		//fmt.Println("- error")
		wg.Done()
	}
	fmt.Println("all done")
	//wg.Done()
	/*for err := range chError {
		fmt.Printf("Eror text %v \n", err)
		fmt.Println("- after")
		wg.Done()
	}*/

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
	//var wgError sync.WaitGroup
	//wgError.Add(1)
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
		//fmt.Println("+")
		select {
		case <-quit:
			//fmt.Println("stop add job")
			//close(workerInput)
			wg.Done()
			//fmt.Println("-")
			//workerInput <- f
			//break
		default:
			fmt.Println("add job")
			workerInput <- f
		}
		//wg.Done()
	}
	close(workerInput)
	//close(chError)
	wg.Wait()
	fmt.Println("end wait")

	/*for err := range chError {
		fmt.Printf("Eror text %v \n", err)
	}*/

}

func work(x int) func() error {
	return func() error {
		/*for i := 1; i > 0; i++ {
			if i%x*1000 == 0 {
				time.Sleep(1000 * time.Millisecond)
				return fmt.Errorf("fail work: %d", x)
			}
		}*/
		return nil
	}
}
