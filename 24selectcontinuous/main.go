package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const bufSize = 5

//複数のChannelの値を連続して受け取る方法

func main() {
	ch1 := make(chan int, bufSize)
	ch2 := make(chan int, bufSize)
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Millisecond)
	defer cancel()
	wg.Add(3)
	go countProducer(&wg, ch1, bufSize, 50)
	go countProducer(&wg, ch2, bufSize, 500)
	go countConsumer(ctx, &wg, ch1, ch2)
	wg.Wait()
	fmt.Println("finish")

}

// Channelへ書き込み用の関数
func countProducer(wg *sync.WaitGroup, ch chan<- int, size int, sleep int) {
	defer wg.Done()
	defer close(ch)
	for i := 0; i < size; i++ {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		ch <- i
	}
}

// Channelから受信を行うための関数
func countConsumer(ctx context.Context, wg *sync.WaitGroup, ch1 <-chan int, ch2 <-chan int) {
	defer wg.Done()
	//loop:
	for ch1 != nil || ch2 != nil {
		select {
		case <-ctx.Done():
			//タイムアウトした場合の動作。
			//タイムアウトした場合、ch1, ch2はnilになる。
			//このため、select文のcaseには入らない。
			//Done()は、contextのメソッド。構造体のChannelを返す。
			fmt.Println(ctx.Err())
			//break loop

			//タイムアウト後の処理をここに書くことで、Channelのバッファに残っていた処理を実行できる。
			//fotとselectが二重構造になるのが、冗長であまり美しくない気もするが。
			for ch1 != nil || ch2 != nil {
				select {
				case v, ok := <-ch1:
					if !ok {
						//okにfalseが返ってきた場合は、Channelが閉じられている。
						ch1 = nil
						break
					}
					fmt.Printf("ch1: %v\n", v)
				case v, ok := <-ch2:
					if !ok {
						ch2 = nil
						break
					}
					fmt.Printf("ch2: %v\n", v)
				}
			}
		case v, ok := <-ch1:
			if !ok {
				//okにfalseが返ってきた場合は、Channelが閉じられている。
				ch1 = nil
				break
			}
			fmt.Printf("ch1: %v\n", v)
		case v, ok := <-ch2:
			if !ok {
				ch2 = nil
				break
			}
			fmt.Printf("ch2: %v\n", v)
		}

	}

}
