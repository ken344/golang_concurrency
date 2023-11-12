package main

import "fmt"

// select文は、複数のチャネルから受信することができる.
// どのチャネルから受信するかは、ランダムに決定される.
func main() {
	ch1 := make(chan string, 1)
	ch2 := make(chan string, 1)

	ch1 <- "ch1"
	ch2 <- "ch2"

	select {
	case v := <-ch1:
		fmt.Println(v)
	case v := <-ch2:
		fmt.Println(v)
	}
}
