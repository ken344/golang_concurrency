package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
	"time"
)

func main() {
	//var wg sync.WaitGroup
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	fmt.Println("goroutine invoked")
	//}()
	//wg.Wait()
	//fmt.Printf("run of working groutines: %d\n\n", runtime.NumGoroutine())
	//fmt.Println("main func finished")

	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal("Error:", err)
		}
	}()
	if err := trace.Start(f); err != nil {
		log.Fatalln("Error:", err)
	}
	defer trace.Stop()
	ctx, t := trace.NewTask(context.Background(), "main")
	defer t.End()
	fmt.Println("The number of logical CPU Cores", runtime.NumCPU())

	//task(ctx, "task1")
	//task(ctx, "task2")
	//task(ctx, "task3")
	var wg sync.WaitGroup
	wg.Add(3)
	go cTask(ctx, &wg, "task1")
	go cTask(ctx, &wg, "task2")
	go cTask(ctx, &wg, "task3")
	wg.Wait()

	s := []int{1, 2, 3}
	for _, i := range s {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}
	wg.Wait()

	fmt.Println("main func finished")
}

// 通常の実行版
func task(ctx context.Context, name string) {
	defer trace.StartRegion(ctx, "task").End()
	time.Sleep(time.Second)
	fmt.Println(name)
}

// 並列実行版
func cTask(ctx context.Context, wg *sync.WaitGroup, name string) {
	defer trace.StartRegion(ctx, "task").End()
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Println(name)
}
