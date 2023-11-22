package main

import (
	"fmt"
	"time"
)

const bufferSize = 5
const numOfReaders = 2
const numOfWriters = 8

func main() {
	// Shared channel between producer and consumer
	dataChannel := make(chan int, bufferSize)
	done := make(chan struct{}) // Signal to terminate the consumer

	// Start the producer
	for i := 1; i < numOfWriters; i++ {
		go producer(dataChannel)
	}
	// Start the consumer

	for i := 1; i < numOfReaders; i++ {
		go consumer(dataChannel, done)
	}

	// Let the program run for a while
	time.Sleep(time.Second * 5)

	// Signal the producer to stop
	close(done)

	// Wait for the consumer to finish
	time.Sleep(time.Second)
}

func producer(dataChannel chan<- int) {

	for i := 1; i <= 10; i++ {

		time.Sleep(time.Millisecond * 500)
		// Produce data (an integer in this case)
		fmt.Println("Producing:", i)
		dataChannel <- i
	}

	fmt.Println("Producer finished producing.")
}

func consumer(dataChannel <-chan int, done <-chan struct{}) {
	for { //go routines are assumed to always be running, an infinite loop

		data, ok := <-dataChannel
		if !ok {
			fmt.Println("Consumer finished consuming.")
		} else {
			time.Sleep(time.Millisecond * 200)
			// Consume data (print it in this case)
			fmt.Println("Consuming:", data)
		}

	}
}
