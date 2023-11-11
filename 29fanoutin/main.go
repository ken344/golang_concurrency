package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	cores := runtime.NumCPU()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}

	// coresの数だけ、chanのスライスを作成する.
	outChs := make([]<-chan string, cores)
	//可変長引数を渡す場合は、スライスを渡す.
	inData := generator(ctx, nums...)
	for i := 0; i < cores; i++ {
		// fanOut関数は、inチャネルから読み取った値を、outチャネルのスライスに書き込む.
		outChs[i] = fanOut(ctx, inData, i+1)
	}
	var i int
	flag := true
	for v := range fanIn(ctx, outChs...) {
		if i == 3 {
			cancel()
			flag = false
		}
		if flag {
			fmt.Println(v)
		}
		i++
	}
	fmt.Println("finish")
}

// 可変長で渡した数値を、分割し、チャネルに書き込む.
func generator(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case <-ctx.Done():
				return
			case out <- n:
			}
		}

	}()
	return out
}

func fanOut(ctx context.Context, in <-chan int, id int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		// 重い処理をシミュレートする.
		heavyWork := func(i int, id int) string {
			time.Sleep(200 * time.Millisecond)
			return fmt.Sprintf("result:%v (id:%v)", i*i, id)
		}
		// inチャネルに書き込みがある度に、チャネルから読み取った値を、heavyWork関数に渡して、outチャネルに書き込む.
		for v := range in {
			select {
			case <-ctx.Done():
				return
			case out <- heavyWork(v, id):
			}
		}

	}()
	return out
}

// fanOutで複数のチャネルに分割した処理を、fanInで一つのチャネルにまとめる.
func fanIn(ctx context.Context, chs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)
	//チャネルに書き込みがある度に、for rangeで取り出し、出力Channelに書き込む
	multiplex := func(ch <-chan string) {
		defer wg.Done()
		for text := range ch {
			select {
			case <-ctx.Done():
				return
			case out <- text:
			}
		}
	}
	//可変長のチャネルの数分だけ、Goroutineを生成し、multiplex関数を実行する.
	wg.Add(len(chs))
	for _, ch := range chs {
		//引数は、chの値を渡す
		go multiplex(ch)
	}
	//立ち上げたすべてのGoroutineが終了するまで待つ.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
