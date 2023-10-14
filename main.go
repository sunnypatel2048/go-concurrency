package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	lock  sync.Mutex
	count int
)

func main() {
	// evilNinjas := []string{"Tommy", "Jony", "Bobby", "Andy"}
	// withoutConcurrency(evilNinjas)
	// withConcurrencyNoChannel(evilNinjas)
	// withConcurrencyWithChannel("Tommy")
	// withBufferedChannel()
	// withConcurrencyWithChannelIterationAndClosing()
	// withConcurrencyWithWaitGroup(evilNinjas)
	// withConcurrencyWithMutex()
}

func withoutConcurrency(evilNinjas []string) {
	start := time.Now()
	defer func() {
		fmt.Println("Attack without concurrency took", time.Since(start))
	}()

	for _, evilNinja := range evilNinjas {
		attack(evilNinja)
	}
}

func withConcurrencyNoChannel(evilNinjas []string) {
	start := time.Now()
	defer func() {
		fmt.Println("Attack with concurrency (No Channels) took", time.Since(start))
	}()

	for _, evilNinja := range evilNinjas {
		go attack(evilNinja)
	}

	// Need to add some sort of delay
	// otherwise other go routines will not get time to execute before main go routine stops.
	time.Sleep(time.Second * 2)
}

func withConcurrencyWithChannel(evilNinja string) {
	start := time.Now()
	defer func() {
		fmt.Println("Attack with concurrency (using channel) took", time.Since(start))
	}()

	// channel is like shared memory, which can be accessed by main go routine and other go routines
	smokeSignal := make(chan bool)
	go attackWithChannel(evilNinja, smokeSignal)
	fmt.Println(<-smokeSignal)
}

func withBufferedChannel() {
	// By default channel capacity is 0, So message needs to be unloaded as soon as it is added to the channel
	// This execution will throw deadlock errror, if there is nowhere the message in the channel can be unloaded to.
	// To avoid this we can create channel with specified capacity.
	channel := make(chan string, 1)
	channel <- "First message"
	// We have a channel with capacity 1, if below line is uncommented code will throw deadlock error
	// channel <- "Second message"
	fmt.Println(<-channel)
}

func withConcurrencyWithChannelIterationAndClosing() {
	channel := make(chan string)
	go throwNinjaStar(channel)
	for {
		message, open := <-channel
		if !open {
			break
		}
		fmt.Println(message)
	}
}

func withConcurrencyWithWaitGroup(evilNinjas []string) {
	var beeper sync.WaitGroup
	beeper.Add(len(evilNinjas))
	for _, evilNinja := range evilNinjas {
		go attackWithWaitGroup(evilNinja, &beeper)
	}
	beeper.Wait()
	fmt.Println("Mission Completed!")
}

func withConcurrencyWithMutex() {
	iterations := 1000
	for i := 0; i < iterations; i++ {
		go increment()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Result Count is:", count)
}

func attack(target string) {
	fmt.Println("Throwing ninja stars at", target)
	time.Sleep(time.Second)
}

func attackWithChannel(target string, attacked chan bool) {
	time.Sleep(time.Second)
	fmt.Println("Throwing ninja stars at", target)
	attacked <- true
}

func attackWithWaitGroup(evilNinja string, beeper *sync.WaitGroup) {
	fmt.Println("Attacked evil ninja: ", evilNinja)
	beeper.Done()
}

func throwNinjaStar(channel chan string) {
	rand.NewSource(time.Now().UnixNano())
	numRounds := 3
	for i := 0; i < numRounds; i++ {
		score := rand.Intn(10)
		channel <- fmt.Sprint("You scored: ", score)
	}
	close(channel)
}

func increment() {
	lock.Lock()
	count++
	lock.Unlock()
}
