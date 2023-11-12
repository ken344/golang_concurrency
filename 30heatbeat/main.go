package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// heartbeatパターンは、goroutineが正常に動作しているかを確認するために、
// 定期的に信号を送るパターン.
// ここでは、task関数が正常に動作しているかを確認するために、
// heartbeatパターンを使用している.
func main() {
	// ログファイルを作成する.
	file, err := os.Create("log.txt")
	if err != nil {
		log.Fatalln(err)
	}
	// ログファイルを閉じる.
	defer file.Close()
	// ログファイルと標準エラー出力に出力する.
	// ログのプレフィックスをERROR:にする.
	// ログの発生時刻を出力する.
	errorLogger := log.New(io.MultiWriter(file, os.Stderr), "ERROR: ", log.LstdFlags)
	//タイムアウトを設定する
	ctx, cancel := context.WithTimeout(context.Background(), 5100*time.Millisecond)
	defer cancel()
	//ウォッチドックタイマーとハートビートタイマーを定義する
	const wdtTimeout = 800 * time.Millisecond
	const beatInterval = 500 * time.Millisecond
	heartbeat, v := task(ctx, beatInterval)
loop:
	//waatchdogタイマーを定義する
	//for文になっているので、select文のcaseのどれかが実行され続ける
	for {
		select {
		//ハートビートタイマーを受信したときに、メッセージを出力
		case _, ok := <-heartbeat:
			if !ok {
				break loop
			}
			fmt.Println("beat pulse")
		//一秒周期で送られてくるValueの値がチャネルに書き込まれている場合に実行
		case r, ok := <-v:
			if !ok {
				break loop
			}
			//時刻の値は末尾にm=が付与されているので、それを削除して出力（モノトニッククロック）
			t := strings.Split(r.String(), "m=")
			fmt.Printf("value: %v [s]\n", t[1])
		//watchdogタイマーがタイムアウトした場合に実行
		//ハートビートタイマーが到達し続けるかぎり、3つめのcaseは実行されない
		case <-time.After(wdtTimeout):
			errorLogger.Println("doTask goroutine`s heartbeat stopped")
			break loop

		}
	}

}

// task関数は、定期的に値を送信する.
// heartbeatパターンを使用して、task関数が正常に動作しているかを確認する.
// 戻り値は、heartbeatに対応していて、空の読み取り専用のstructを送信するチャネルと、
// task関数が値を送信するチャネル.
func task(ctx context.Context, beatInterval time.Duration) (<-chan struct{}, <-chan time.Time) {
	//ハートビートタイマーは値を持たない通知用のチャネルなので、空の構造体を定義
	heartbeat := make(chan struct{})
	//値を送信するチャネルを定義
	//今回は、tickerを使用して、定期的に値を送信する.
	out := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(out)
		pulse := time.NewTicker(beatInterval)
		task := time.NewTicker(2 * beatInterval)
		//ハートビートの値をチャネルに書き込む関数
		sendPulse := func() {
			select {
			//ハートビートが書き込み可能な場合は、空の構造体を書き込む.
			case heartbeat <- struct{}{}:
			//書き込み可能でない場合は、何もせずにselect文を抜ける.
			default:
			}
		}
		//値をチャネルに書き込む為の関数
		sendValue := func(t time.Time) {
			for {
				select {
				case <-ctx.Done():
					return
				//ハートビートがチャネルに書き込まれたら、sendPulseを実行し、
				//再度for文を実行する
				case <-pulse.C:
					sendPulse()
				//受け取った値を出力Channelに書き込む
				case out <- t:
					return
				}
			}
		}
		//定義した無名関数を呼び出すselect文
		for {
			select {
			case <-ctx.Done():
				return
			case <-pulse.C:
				sendPulse()
			case t := <-task.C:
				sendValue(t)
			}
		}
	}()
	return heartbeat, out
}
