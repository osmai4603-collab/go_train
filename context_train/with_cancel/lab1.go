package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

func StartGoroutine(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	fmt.Println("WORKER: Start working in lab1:")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("WORKER: we receive stop working from main context")
			fmt.Printf("WORKER: cause of stop working: %s\n", ctx.Err())
			return
		case duration := <-ticker.C:
			fmt.Printf("WORKER: Stil Working... %s\n", strconv.Itoa(duration.Second()))
		}
	}

}

func Lab1() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	fmt.Println("======================================== Start lab1 ========================================")
	wg.Add(1)
	go StartGoroutine(ctx, wg)

	time.Sleep(5 * time.Second)

	fmt.Println("BOSS: We will stop worker of working goroutine...")
	cancel()
	wg.Wait()

	fmt.Println("======================================== End lab1 ========================================")

}
