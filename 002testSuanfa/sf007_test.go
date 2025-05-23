package _02testSuanfa

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"
)

// goroutine 阻塞示例

func Test007Name(t *testing.T) {
	runtime.GOMAXPROCS(1) // 强制单核
	go func() {
		for {
			time.Sleep(time.Second) //如果不加这个 goroutine会一直占用CPU，导致主协程得不到调度
		}
	}()
	time.Sleep(time.Second * 5)
	println("主协程被调度") // 观察何时输出

	fmt.Println("GODEBUG:", os.Getenv("GODEBUG"))
}

//好的，用户现在的问题是“GODEBUG=asyncpreemptoff=1 怎么修改”。我需要仔细分析这个问题，确保正确理解用户的需求，并提供有效的解决方案。
//
//首先，用户之前的问题是关于在Go 1.14及之后版本中如何实现无限循环，特别是当使用runtime.GOMAXPROCS(1)时，主协程能否被调度。之前的对话中，我解释了Go 1.14引入的抢占式调度机制，并提到可以通过设置GODEBUG=asyncpreemptoff=1来禁用异步抢占，以测试旧版本的行为。现在用户询问如何具体修改这个环境变量，可能是在实际使用中遇到了问题，或者想要测试禁用抢占后的效果。
//
//接下来，我需要确定用户的具体需求。用户可能想知道如何在不同的环境中设置这个环境变量，比如在命令行中运行程序时，或者在使用IDE（如Goland、VSCode）调试时，或者是在Docker容器中运行。此外，用户可能希望了解这个设置的具体影响，以及如何验证是否生效。
//
//然后，我需要考虑不同的操作系统环境。比如，在Linux/macOS和Windows下，设置环境变量的方式有所不同。用户可能不熟悉如何在Windows中设置，所以需要详细说明。同时，用户可能在使用Go模块或其他工具，需要确认设置是否会影响整个系统或仅当前会话。
//
//另外，用户可能需要知道如何验证设置是否成功。比如，在代码中加入调试输出，或者观察程序行为的变化。例如，当禁用异步抢占后，原来的无限循环goroutine是否会阻止主协程执行，从而验证环境变量是否生效。
//
//还可能存在的潜在问题是，用户可能误写环境变量名称或值，导致设置不生效。需要提醒常见的错误，比如拼写错误，或者是否在正确的shell中设置。此外，需要说明这个设置仅适用于Go 1.14及以上版本，旧版本可能不支持。
//
//另外，用户可能想知道这个环境变量的作用时间，是否是临时的还是永久的。需要解释在命令行中设置是临时的，仅影响当前进程，而修改系统环境变量则是永久性的，但通常建议仅在需要时临时设置。
//
//最后，需要总结使用场景，比如调试、测试旧版本行为，或者性能分析，同时提醒用户在生产环境中不应禁用此功能，以避免导致协程饿死等问题。确保回答清晰、全面，覆盖用户可能的各种使用场景和疑问。
//
//
//
//
//在 Go 1.14 及更高版本中，如果你想临时禁用基于信号的异步抢占式调度（例如用于调试或测试旧版本行为），可以通过设置 `GODEBUG` 环境变量来实现。以下是具体操作步骤：
//
//---
//
//### **1. 通过命令行临时设置（推荐）**
//在运行 Go 程序时，直接在命令行中设置 `GODEBUG` 环境变量：
//
//#### **Linux/macOS**
//```bash
//GODEBUG=asyncpreemptoff=1 go run main.go
//```
//
//#### **Windows（PowerShell）**
//```powershell
//$env:GODEBUG = "asyncpreemptoff=1"
//go run main.go
//```
//
//#### **Windows（CMD）**
//```cmd
//set GODEBUG=asyncpreemptoff=1
//go run main.go
//```
//
//---
//
//### **2. 在 IDE 中设置（如 Goland、VSCode）**
//如果你使用 IDE 调试程序，可以在运行配置中添加 `GODEBUG` 环境变量：
//
//#### **Goland**
//1. 打开 `Run/Debug Configurations`。
//2. 在 `Environment` 字段中添加：
//   ```
//   GODEBUG=asyncpreemptoff=1
//   ```
//
//#### **VSCode**
//在 `.vscode/launch.json` 中添加环境变量：
//```json
//{
//    "version": "0.2.0",
//    "configurations": [
//        {
//            "name": "Launch",
//            "type": "go",
//            "request": "launch",
//            "env": {
//                "GODEBUG": "asyncpreemptoff=1"
//            }
//        }
//    ]
//}
//```
//
//---
//
//### **3. 在 Docker 容器中设置**
//在 Dockerfile 或 `docker run` 命令中指定环境变量：
//```dockerfile
//ENV GODEBUG=asyncpreemptoff=1
//```
//或
//```bash
//docker run -e GODEBUG=asyncpreemptoff=1 my-go-app
//```
//
//---
//
//### **4. 验证设置是否生效**
//可以在代码中添加调试输出，观察程序行为是否符合预期：
//```go
//package main
//
//import (
//	"fmt"
//	"os"
//	"runtime"
//	"time"
//)
//
//func main() {
//	runtime.GOMAXPROCS(1)
//	go func() {
//		for {} // 无限循环
//	}()
//	time.Sleep(time.Second)
//	fmt.Println("主协程被调度") // 观察是否输出
//
//	// 打印环境变量（可选）
//	fmt.Println("GODEBUG:", os.Getenv("GODEBUG"))
//}
//```
//
//#### **禁用抢占时的行为**：
//• 如果设置 `asyncpreemptoff=1`，主协程**不会输出** `主协程被调度`（无限循环协程饿死主协程）。
//• 如果未设置或设置错误，主协程会正常输出（抢占生效）。
//
//---
//
//### **注意事项**
//1. **临时性**：
//   通过命令行或 IDE 设置的环境变量仅对当前进程生效，重启程序后失效。
//
//2. **兼容性**：
//   `asyncpreemptoff=1` 仅对 Go 1.14+ 有效。旧版本 Go 无此功能。
//
//3. **生产环境**：
//   **不要在生产环境中禁用抢占式调度**！这会导致协程饿死（如你的示例代码中的主协程无法执行）。
//
//4. **调试用途**：
//   该选项主要用于调试抢占式调度相关的问题，例如对比旧版本和新版本的行为差异。
//
//---
//
//通过以上方法，你可以灵活控制 Go 调度器的行为，验证代码在不同场景下的表现。
