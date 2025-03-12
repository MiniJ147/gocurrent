/*
Benchmark file to test queue vs channel throughput.
Args - Simulation time (seconds), Num Readers, Num Writers

# Outputs total added, removed, and difference between Queue and Channel

Runs number of readers and number of writers on a shared queue/channel and runs for number of seconds
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/minij147/gocurrent/queue"
)

var TIME = 10 * time.Second
var WRITERS = 1
var READERS = 1

// added, killed
func queueTest() (int64, int64) {
	h := queue.New[int]()

	ctx, cancel := context.WithCancel(context.Background())
	var added, killed int64 = 0, 0

	for range WRITERS {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					h.Push(0)
					added++
				}
			}
		}(ctx)
	}

	for range READERS {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					_, success := h.Pop()
					if success {
						killed++
					}
				}
			}
		}(ctx)
	}

	time.Sleep(TIME)
	cancel()
	log.Println("queue test done")
	return added, killed
}

// writes, deletes
func channelTest() (int64, int64) {
	c := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	var added, killed int64 = 0, 0

	for range WRITERS {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					c <- 0
					added++
				}
			}
		}(ctx)
	}

	for range READERS {
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					<-c
					killed++
				}
			}
		}(ctx)
	}

	time.Sleep(TIME)
	cancel()
	log.Println("channel test done")
	return added, killed
}

func main() {
	if len(os.Args) < 4 {
		log.Fatal("not enough args, SIM-TIME seconds, READERS, WRITERS")
		return
	}

	argSimTime, argReads, argWrites := os.Args[1], os.Args[2], os.Args[3]

	simTime, err := strconv.Atoi(argSimTime)
	if err != nil {
		panic(err)
	}
	TIME = time.Duration(simTime) * time.Second

	READERS, err = strconv.Atoi(argReads)
	if err != nil {
		panic(err)
	}

	WRITERS, err = strconv.Atoi(argWrites)
	if err != nil {
		panic(err)
	}

	fmt.Printf("====================\n")
    fmt.Printf("Queue Test starting time: %v\nTime: %vs\nReaders: %v\nWriters: %v\n\n", time.Now(),simTime, READERS, WRITERS)

	log.Println("starting queue test")
	qa, qk := queueTest()

	log.Println("ended queue test\nstarting channel test\n")
	ca, ck := channelTest()

	fmt.Printf("queue\n%v %v %v\n\n", qa, qk, qa-qk)
	fmt.Printf("channel\n%v %v %v\n\n", ca, ck, ca-ck)

	fmt.Printf("Push Diff: %v x%0.2f more | Pop Diff: %v x%0.2f more\n", qa-ca, float64(qa)/float64(ca), qk-ck, float64(qk)/float64(ck))

}
