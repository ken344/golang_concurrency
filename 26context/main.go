package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	//withTimeoutは、指定した時間が経過すると、Contextをキャンセルする.
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
