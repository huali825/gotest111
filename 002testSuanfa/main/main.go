package main

import (
	"flag"
	"fmt"
	"os"
)

//输出内容解析
//当设置该参数后，程序运行时输出的每一行信息包含以下字段（示例）：
//
//sched 1003ms: gomaxprocs=4 idleprocs=0 threads=5 spinningthreads=0 idlethreads=0 runqueue=14 [1 0 1 0]
//各字段含义如下：
//
//​sched 1003ms：程序启动后到当前输出行的时间（1003 毫秒）。
//​gomaxprocs：当前逻辑处理器（P）的总数，通常等于 GOMAXPROCS 的设置值或 CPU 核心数。
//​idleprocs：空闲状态的 P 数量（未绑定 Goroutine 的处理器）。
//​threads：当前运行时管理的操作系统线程（M）总数。
//​spinningthreads：处于“自旋”状态的线程数（正在寻找可运行 Goroutine 的线程）。
//​idlethreads：空闲状态的线程数（未绑定 P 的线程）。
//​runqueue：全局运行队列（GRQ）中等待的 Goroutine 数量。
//​**[1 0 1 0]**：各 P 的本地运行队列（LRQ）中的 Goroutine 数量（例如，4 个 P 的队列分别有 1、0、1、0 个 G）。

// powershell
// $env:GODEBUG = "schedtrace=1000"; go run main.go

//func main() {
//	ch := make(chan int, 2)
//	ch <- 1
//	ch <- 2
//	go func() {
//		time.Sleep(2 * time.Second)
//		<-ch
//	}()
//	fmt.Println("开始阻塞发送...")
//	ch <- 3 // 观察goroutine阻塞状态（GODEBUG=schedtrace=1000）
//	fmt.Println("发送完成")
//}

func main() {
	// 使用flag包定义一个命令行参数，参数名为"n"，默认值为"world"，参数说明为"a string"
	name := flag.String("n", "world", "a string")

	// 解析命令行参数
	flag.Parse()

	// 打印操作系统传递给程序的参数列表
	fmt.Println("os arg is:", os.Args)

	// 打印用户输入的命令行参数值
	fmt.Println("input parameter is :", *name)

	// 使用Sprintf格式化字符串，将"hello "和用户输入的参数值拼接成一个新的字符串
	fullString := fmt.Sprintf("hello %s", *name)

	// 打印拼接后的字符串
	fmt.Println(fullString)
}
