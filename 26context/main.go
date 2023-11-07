package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	// withDeadline()は、指定した時間が経過すると、一斉にサブGoroutineのContextをキャンセルする.
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(20*time.Millisecond))
	defer cancel()
	ch := subTask2(ctx)
	v, ok := <-ch
	if ok {
		fmt.Println(v)
	}
	fmt.Println("finish")
}

func subTask2(ctx context.Context) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		deadline, ok := ctx.Deadline()
		if ok {
			// deadline.Sub(time.Now())は、指定した時間と、現在時刻の差分を返す.
			// 今回は、deadline.Sub(time.Now())が、30ミリ秒より小さい場合は、
			// タイムアウトとして処理を終了する.
			if deadline.Sub(time.Now().Add(30*time.Millisecond)) < 0 {
				fmt.Println("impossible to meet deadline")
				return
			}
		}
		time.Sleep(30 * time.Millisecond)
		ch <- "hello"
	}()
	return ch
}

func main_withCancel() {
	var wg sync.WaitGroup
	// Contextは、親子関係を持つことができる.
	// context.Background()は、親Contextを持たないContextを生成する.
	//withCancelは、キャンセルされた場合に、一斉にサブGoroutineのContextをキャンセルする.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := criticalTask(ctx)
		if err != nil {
			fmt.Printf("critical task caancelled due to: %v\n", err)
			// ここで、cancelを呼び出すことで、
			// criticalTaskのサブGoroutineのContextもキャンセルされる.
			cancel()
			return
		}
		fmt.Println("success", v)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		v, err := normalTask(ctx)
		if err != nil {
			fmt.Printf("normal task cancelled due to: %v\n", err)
			return
		}
		fmt.Println("success", v)
	}()
	wg.Wait()
}

func criticalTask(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 800*time.Millisecond)
	defer cancel()
	t := time.NewTicker(1000 * time.Millisecond)
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-t.C:
		t.Stop()
	}
	return "A", nil
}

func normalTask(ctx context.Context) (string, error) {
	t := time.NewTicker(3000 * time.Millisecond)
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-t.C:
		t.Stop()
	}
	return "B", nil
}

func main_withtimeout() {
	var wg sync.WaitGroup
	// Contextは、親子関係を持つことができる.
	// context.Background()は、親Contextを持たないContextを生成する.
	//withTimeoutは、指定した時間が経過すると、一斉にサブGoroutineのContextをキャンセルする.
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()
	wg.Add(3)
	go subTask(ctx, &wg, "a")
	go subTask(ctx, &wg, "b")
	go subTask(ctx, &wg, "c")
	wg.Wait()
}

func subTask(ctx context.Context, wg *sync.WaitGroup, id string) {
	defer wg.Done()
	// NewTickierは、指定した時間ごとに、チャネルに値を送信する.
	// Channelへの書き込みをシミュレートする。
	t := time.NewTicker(500 * time.Millisecond)
	select {
	case <-ctx.Done():
		//ContextのDoneメソッドは、Contextがキャンセルされた場合に（親からのタイムアウトなどで）、
		//チャネルに値を送信する.
		fmt.Println(ctx.Err())
		return
	case <-t.C:
		//nwTickerのチャネルから値を受信するには、t.Cを利用する.
		//今回は、最初の受信をした後に、newTickerをstopで停止するようにしている
		t.Stop()
		fmt.Println(id)
		return
	}
}
