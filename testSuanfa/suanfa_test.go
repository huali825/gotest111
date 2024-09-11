package testSuanfa

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	printFunc2()
	printFunc2()
}

func printFunc1() {
	func111 := bibaoFunc()
	fmt.Println(func111(1))
	fmt.Println(func111(1))
	//这个函数打印 1, 2  在这个函数变量 func111 定义之后, sum 的修改会影响func1全局
}

func printFunc2() {
	func111 := bibaoFunc()
	fmt.Println(func111(1))
	fmt.Println(func111(1))
}

func bibaoFunc() func(val int) int {
	sum := 0
	return func(val int) int {
		sum += val
		return sum
	}
}
