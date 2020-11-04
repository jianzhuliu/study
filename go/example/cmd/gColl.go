/*
gColl.go 

linux:	 
	GODEBUG=gctrace=1 go run gColl.go
	
windows:
	set GODEBUG=gctrace=1 
	go run gColl.go 

Output:
gc 1 @0.080s 0%: 0+1.0+0 ms clock, 0+0/0/0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 2 @0.103s 0%: 1.0+0+0 ms clock, 4.0+0/0/0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 3 @0.142s 0%: 0+0+0 ms clock, 0+0/0/0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 4 @0.183s 0%: 0+1.0+0 ms clock, 0+1.0/1.0/2.0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 5 @0.229s 0%: 0+0.97+0 ms clock, 0+0/0.97/0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 6 @0.415s 0%: 0+1.0+0 ms clock, 0+0/1.0/1.0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 7 @0.440s 0%: 0+0.99+0 ms clock, 0+0/0.99/1.9+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 8 @0.461s 0%: 1.0+1.0+0 ms clock, 4.0+0/1.0/1.0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
gc 9 @0.471s 0%: 0+2.0+0 ms clock, 0+2.0/2.0/2.0+0 ms cpu, 4->4->0 MB, 5 MB goal, 4 P
# command-line-arguments
gc 1 @0.009s 4%: 0+2.0+0 ms clock, 0+1.0/1.0/0+0 ms cpu, 4->6->5 MB, 5 MB goal, 4 P
gc 2 @0.014s 6%: 0+3.0+0 ms clock, 0+0/3.0/0+0 ms cpu, 12->12->12 MB, 13 MB goal, 4 P
gc 3 @0.031s 6%: 0.99+2.9+0 ms clock, 3.9+0/1.0/1.9+0 ms cpu, 21->22->19 MB, 24 MB goal, 4 P
gc 4 @0.098s 3%: 0+8.0+0 ms clock, 0+0/6.9/2.0+0 ms cpu, 36->37->25 MB, 39 MB goal, 4 P
mem.Alloc: 112672
mem.TotalAlloc: 112672
mem.HeapAlloc: 112672
mem.NumGC: 0
----------------
gc 1 @0.010s 0%: 0+1.0+0 ms clock, 0+0/0/2.0+0 ms cpu, 4->5->1 MB, 5 MB goal, 4 P
mem.Alloc: 1116504
mem.TotalAlloc: 5194720
mem.HeapAlloc: 1116504
mem.NumGC: 1
----------------	
	
	
这些数据提供了更多垃圾回收过程中的堆内存大小的信息，
gc 3 @0.031s 6%: 0.99+2.9+0 ms clock, 3.9+0/1.0/1.9+0 ms cpu, 21->22->19 MB, 24 MB goal, 4 P
gc 3 	表示第3次执行
@0.031s 	表示程序执行的总时间
6%: 	垃圾回收时间占用的百分比
0.99+2.9+0 ms clock 	垃圾回收时间，分别为STW清扫时间，并发标记和扫描时间，STW标记时间
3.9+0/1.0/1.9+0 ms cpu	垃圾回收占用 cpu 时间
21->22->19 MB 垃圾回收器要去运行时候的堆内存大小，垃圾回收器操作结束时候的堆内存大小，存活堆的大小
24 MB goal	整体堆的大小
4 P	使用的处理器数量

	
*/
package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
)

func showErr(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(2)
}

func printMem(mem runtime.MemStats) {
	runtime.ReadMemStats(&mem)
	fmt.Println("mem.Alloc:", mem.Alloc)
	fmt.Println("mem.TotalAlloc:", mem.TotalAlloc)
	fmt.Println("mem.HeapAlloc:", mem.HeapAlloc)
	fmt.Println("mem.NumGC:", mem.NumGC)
	fmt.Println("----------------")
}

func main() {
	f, err := os.Create("traceFile.out")
	if err != nil {
		showErr("os.Open", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			showErr("f.Close()", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		showErr("trace.Start()", err)
	}
	defer trace.Stop()

	var mem runtime.MemStats
	printMem(mem)

	for i := 0; i < 10; i++ {
		s := make([]byte, 500000)
		if s == nil {
			fmt.Println("Operation failed")
		}
	}

	printMem(mem)
}
