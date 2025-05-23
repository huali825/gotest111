package _03goroutineCode

import (
	"fmt"
	"testing"
	"time"
)

type printNumAbc struct {
	numValue  int
	byteValue byte
	ch1       chan int
	ch2       chan int
}

func (p *printNumAbc) printNum() {

	for i := 0; i < 14; i++ {
		<-p.ch1
		fmt.Print(p.numValue)
		p.numValue++
		fmt.Print(p.numValue)
		p.numValue++
		p.ch2 <- 0
	}

}

func (p *printNumAbc) printAbc() {

	for i := 0; i < 13; i++ {
		<-p.ch2
		fmt.Printf("%c", p.byteValue)
		p.byteValue++
		fmt.Printf("%c", p.byteValue)
		p.byteValue++
		p.ch1 <- 0
	}

}

func TestPrintNumAbc(t *testing.T) {
	p := printNumAbc{1, 'A', make(chan int), make(chan int)}
	go p.printNum()
	go p.printAbc()
	p.ch1 <- 0

	//fmt.Println("main goroutine over")
	time.Sleep(1 * time.Second)

}
