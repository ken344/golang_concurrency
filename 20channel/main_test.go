package main

import "testing"

import "go.uber.org/goleak"

func TestLeak(t *testing.T) {
	defer goleak.VerifyNone(t)
	main()

}
