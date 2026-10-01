package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

func Generator(ctx context.Context) <-chan int {
	chnl := make(chan int, 5)
	go func() {
		defer close(chnl)
		index := 1
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("🛑 [Generator]: Received stop signal -> %v\n", ctx.Err())
				return
			case chnl <- index:
				index++
			}
		}
	}()
	return chnl
}

func Lab2() {
	fmt.Println("======================================== Start Lab 2 ========================================")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chnl := Generator(ctx)
	fmt.Printf("📊 Goroutines before starting generator: %d\n", runtime.NumGoroutine())

	for range 5 {
		val := <-chnl
		fmt.Printf("📥 [Consumer]: Received number: %d\n", val)
	}

	fmt.Println("\n📢 [Consumer]: Consumed 5 items. Sending cancel() signal...")
	cancel()
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("\n📊 Goroutines after cancellation: %d\n", runtime.NumGoroutine())
	fmt.Println("======================================== End Lab 2 ==========================================")
}
