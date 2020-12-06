package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// SafeCounter is safe to use concurrently.
type SafeCounter struct {
	max   int
	index int
	value []int
	//v   map[string]int
	mux sync.Mutex
}

func new(max int) *SafeCounter {
	n := SafeCounter{}
	n.max = max - 4
	n.value = make([]int, max)
	return &n
}

// Inc increments the counter for the given key.
func (c *SafeCounter) Inc(val int) {
	//if c.index >= c.max {
	//	return
	//}
	c.mux.Lock()
	for c.index >= c.max {
		time.Sleep(time.Microsecond)
	}
	c.value[c.index] = val
	c.index++
	// Lock so only one goroutine at a time can access the map c.v.
	//c.v[key]++
	c.mux.Unlock()
	fmt.Println("added", val)
}

// Value returns the current value of the counter for the given key.
func (c *SafeCounter) Value() (int, error) {
	fmt.Println(c)
	if c.index < 0 {
		return 0, errors.New("empty")
	}
	//c.mux.Lock()
	r := c.value[c.index]
	c.index--
	// Lock so only one goroutine at a time can access the map c.v.
	//c.mux.Unlock()
	return r, nil
}

func run(c *SafeCounter, i int) {
	fmt.Println("go run", i)
	c.Inc(i)
	fmt.Println("stop", i)
}

func main() {
	c := new(6)
	for i := 0; i < 8; i++ {
		//go c.Inc(i)
		go run(c, i)
	}

	time.Sleep(time.Second * 3)
	fmt.Println("end sleep")

	for i := 0; i < 3; i++ {
		r, _ := c.Value()
		fmt.Println("End", r)
	}
	time.Sleep(time.Second * 3)
}
