package _02testSuanfa

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

// 2025年3月11日15:34:18
func Test012Name(t *testing.T) {
	var data [][]byte
	for i := 0; i < 20; i++ {
		// 分配 100MB 内存
		data = append(data, make([]byte, 100*1024*1024))
		time.Sleep(500 * time.Millisecond)
		printMemUsage() // 打印当前内存使用情况
	}
	data = nil
	runtime.GC() // 强制触发 GC
	printMemUsage()
}

// 打印当前内存使用情况
func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Heap: %v MB\n", m.Alloc/1024/1024)
}
