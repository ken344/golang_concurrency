package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var wg sync.WaitGroup
	var c int64

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				//atmicは、排他制御を行うためのライブラリ
				//atomic.AddInt64は、int64型の変数に対して、
				//排他制御を行いながら、値を加算する関数
				//この関数を利用することで、排他制御を行わなくても、
				//複数のgoroutineから、同時に変数に対して、値を加算することができる。
				//第一引数には、加算したい変数のポインタを指定する。
				//第二引数には、加算したい値を指定する。
				atomic.AddInt64(&c, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println(c)
	fmt.Println("finish")
}

func main_RWmutex() {
	// RWmutexは、読み取りと書き込みを分けることができる.
	// Readの場合は、unlockがされなくても、Readを行うことができる.
	// Writeの場合は、unlockがされないと、Writeを行うことができない.
	var wg sync.WaitGroup
	var rwMu sync.RWMutex
	var c int

	wg.Add(4)
	go write(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)
	go read(&rwMu, &wg, &c)

	wg.Wait()
	fmt.Println("finish")
}
func read(mu *sync.RWMutex, wg *sync.WaitGroup, c *int) {
	defer wg.Done()
	time.Sleep(10 * time.Millisecond)
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("read lock")
	fmt.Println(*c)
	time.Sleep(1 * time.Second)
	fmt.Println("read unlock")
}

func write(mu *sync.RWMutex, wg *sync.WaitGroup, c *int) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()
	fmt.Println("write lock")
	*c += 1
	time.Sleep(1 * time.Second)
	fmt.Println("write unlock")
}

func main_mutex() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var i int
	wg.Add(2)
	//2つのgoroutineが同時にiをインクリメントする。
	//このとき、i++の部分は、CPUの命令としては、
	//1. メモリからiの値を読み取る
	//2. CPUのレジスタにiの値を格納する
	//3. レジスタの値をインクリメントする
	//4. レジスタの値をメモリに書き込む
	//という4つの命令に分解される。
	//このとき、goroutineの切り替わりが発生する可能性がある。
	//goroutineの切り替わりが発生すると、メモリからレジスタに値を格納する直前に、
	//別のgoroutineがレジスタの値をインクリメントしてしまう可能性がある。
	//そのため、iの値が2になる可能性がある。
	//また、同時に読み取ることで、最終的に「1」が出力される可能性もある。
	//この現象を、データ競合、もしくは、race conditionと呼ぶ。
	//このような問題を解決するために、mutexを利用する。
	//「go run -race main.go」 で、race conditionが発生しているか確認できる。
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		i++
	}()
	go func() {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		i++
	}()
	wg.Wait()
	fmt.Println(i)
}
