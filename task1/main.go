// Task 1
//Write a program that uses Go's concurrency features to print the numbers from 1 to 100
//concurrently. The program should use a goroutine for each number, and the main goroutine should
//wait for all the goroutines to finish before exiting.

package main

import (
	"fmt"
	"sync"
)

func PrintNumber(num int, wg *sync.WaitGroup) {
	fmt.Println(num)
	wg.Done() // Calling done will decrement the counter

}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 100; i++ {
		wg.Add(1) // incrementing the wait group counter by 1
		go PrintNumber(i, &wg)

	}

	wg.Wait() // main routine will wait for all the go routines to finish up their process
}
