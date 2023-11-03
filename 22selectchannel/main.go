package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)
	var wg sync.WaitGroup
	//タイムアウトを使ったキャンセル機能追加
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	//戻り値として返される関数（Cancelfunc）を呼び出すことで、タイムアウトをキャンセルできる
	//プログラム終了時にCancelFuncを呼び出す。
	defer cancel()

	wg.Add(2)
	go func() {
		defer wg.Done()
		time.Sleep(500 * time.Millisecond)
		ch1 <- "A"
	}()
	go func() {
		defer wg.Done()
		time.Sleep(800 * time.Millisecond)
		ch2 <- "B"
	}()
loop:
	for ch1 != nil || ch2 != nil {
		select {
		case <-ctx.Done():
			//タイムアウトした場合の動作。
			//タイムアウトした場合、ch1, ch2はnilになる。
			//このため、select文のcaseには入らない。
			//Done()は、contextのメソッド。構造体のChannelを返す。
			fmt.Println("timeout")
			break loop
		case v := <-ch1:
			fmt.Println(v)
			ch1 = nil
		case v := <-ch2:
			fmt.Println(v)
			ch2 = nil
		}
	}
	wg.Wait()
	fmt.Println("finished")
}
