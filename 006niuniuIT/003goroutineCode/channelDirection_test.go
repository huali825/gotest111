package _03goroutineCode

import (
	"fmt"
	"testing"
	"time"
)

func TestChannelDirection(t *testing.T) {
	ch := make(chan int, 5)
	go writeChan(ch)
	go readChan(ch)

	time.Sleep(5 * time.Second)
	fmt.Println("main goroutine over")

}

func writeChan(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}

	close(ch)
}

func readChan(ch <-chan int) {
	for value := range ch {
		fmt.Println(value)
	}
}
