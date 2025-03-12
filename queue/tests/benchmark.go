/*
Benchmark file to test queue vs channel throughput.
Args - Simulation time (seconds), Num Readers, Num Writers

# Outputs total added, removed, and difference between Queue and Channel

Runs number of readers and number of writers on a shared queue/channel and runs for number of seconds

WARNING NOT FINISHED, FIX DEADLOCK On Channel Pause
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/minij147/gocurrent/queue"
)

var TIME = 10 * time.Second
var WRITERS = 1
var READERS = 1

// added, killed
func queueTest() (int64, int64) {
	wg := sync.WaitGroup{}
	h := queue.New[int]()

	wg.Add(READERS + WRITERS)
	ctx, cancel := context.WithCancel(context.Background())
	var added, killed int64 = 0, 0

	for range WRITERS {
		go func(ctx context.Context) {
			var localAdded int64 = 0
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					atomic.AddInt64(&added, localAdded)
					return
				default:
					h.Push(0)
					localAdded++
				}
			}
		}(ctx)
	}

	for range READERS {
		go func(ctx context.Context) {
			var localKilled int64 = 0

			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					atomic.AddInt64(&killed, localKilled)
					return
				default:
					_, success := h.Pop()
					if success {
						localKilled++
					}
				}
			}
		}(ctx)
	}

	time.Sleep(TIME)
	cancel()
	wg.Wait()
	log.Println("queue test done")
	return added, killed
}

// writes, deletes
func channelTest() (int64, int64) {
	wg := sync.WaitGroup{}
	wg.Add(READERS + WRITERS)

	c := make(chan int, 10000)
	ctx, cancel := context.WithCancel(context.Background())
	var added, killed int64 = 0, 0

	for range WRITERS {
		go func(ctx context.Context) {
			var localAdded int64 = 0
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					atomic.AddInt64(&added, localAdded)
					return
				default:
					c <- 0
					localAdded++
				}
			}
		}(ctx)
	}

	for range READERS {
		go func(ctx context.Context) {
			var localKilled int64 = 0
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					atomic.AddInt64(&killed, localKilled)
					return
				case <-c:
					localKilled++
				}
			}
		}(ctx)
	}

	time.Sleep(TIME)
	cancel()
	wg.Wait()
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

	log.Printf("====================\n")
	log.Printf("Queue Test starting time: %v\nTime: %vs\nReaders: %v\nWriters: %v\n\n", time.Now(), simTime, READERS, WRITERS)

	log.Println("starting queue test")
	qa, qk := queueTest()

	log.Println("ended queue test\nstarting channel test\n")
	ca, ck := channelTest()
	fmt.Printf("%v,%v,%v,%v,%v\n", simTime, ca, ck, qa, qk)

	log.Printf("queue\n%v %v %v\n\n", qa, qk, qa-qk)
	log.Printf("channel\n%v %v %v\n\n", ca, ck, ca-ck)

	log.Printf("Push Diff: %v x%0.2f more | Pop Diff: %v x%0.2f more\n", qa-ca, float64(qa)/float64(ca), qk-ck, float64(qk)/float64(ck))

}
