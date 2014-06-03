package main

import (
	"fmt"
	"math/rand"
	"sync"
)

var SharedA uint32
var SharedB uint32
var SharedC uint32

func goroutine1(computation <-chan func(), wg *sync.WaitGroup) {
	for {
		f := <-computation
		f()
	}
}

func goroutine2(ch chan func(), wg *sync.WaitGroup) {
	for i := 0; i < 1000000; i++ {
		r := rand.Intn(3)
		ch <- func() {
			switch r {
			case 0:
				SharedA++
			case 1:
				SharedB++
			case 2:
				SharedC++
			}
		}
	}
	wg.Done()

}

func goroutine3(ch <-chan *Output, wg *sync.WaitGroup) {
	for {
		output := <-ch
		switch output.Data {
		case 0:
			SharedA++
		case 1:
			SharedB++
		case 2:
			SharedC++
		}
	}
	wg.Done()
}

type Output struct {
	Data int
}

func goroutine4(ch chan *Output, wg *sync.WaitGroup) {
	for i := 0; i < 1000000; i++ {
		r := rand.Intn(3)
		ch <- &Output{
			Data: r,
		}
	}
	wg.Done()
}

func main() {
	// ch := make(chan func()())
	ch := make(chan *Output)
	wg := &sync.WaitGroup{}
	go goroutine3(ch, wg)

	for i := 0; i < 5; i++ {
		go goroutine4(ch, wg)
		wg.Add(1)
	}

	wg.Wait()
	fmt.Println("main finish")
	return
}
