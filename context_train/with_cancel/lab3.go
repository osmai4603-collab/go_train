package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Parent(parentCTX context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Parent Start working")
	ticker := time.NewTicker(750 * time.Millisecond)
	defer ticker.Stop()
	counter := 1

	for {
		select {
		case <-parentCTX.Done():
			fmt.Printf("🛑 [Parent]: Received stop signal -> %v\n", parentCTX.Err())
			return
		case <-ticker.C:
			fmt.Println("[Parent]: Still Working")
			if counter += 1; counter == 10 {
				cancel()
			}
		}
	}
}

func Client(parentCTX context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(250 * time.Millisecond)

	fmt.Println("Client Start working")
	clientCTX, cancel := context.WithCancel(parentCTX)
	defer cancel()

	ticker := time.NewTicker(750 * time.Millisecond)
	defer ticker.Stop()
	counter := 1
	for {
		select {
		case <-clientCTX.Done():
			fmt.Printf("🛑 [Client]: Received stop signal -> %v\n", clientCTX.Err())
			return

		case <-ticker.C:
			fmt.Println("[Client]: Still Working")
			if counter += 1; counter == 5 {
				cancel()
			}
		}

	}
}

func Notifications(parentCTX context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(500 * time.Millisecond)

	fmt.Println("Notification Start working")
	ctx, notifiCancel := context.WithCancel(parentCTX)
	defer notifiCancel()

	ticker := time.NewTicker(750 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("🛑 [Notifi]: Received stop signal -> %v\n", ctx.Err())
			return
		case <-ticker.C:
			fmt.Println("[Notifi]: Still Working")
		}
	}
}

func Lab3() {
	fmt.Println("======================================== Start Lab 3 ========================================")
	parentCTX, parentCancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(3)

	go Parent(parentCTX, parentCancel, wg)
	go Client(parentCTX, wg)
	go Notifications(parentCTX, wg)

	wg.Wait()

	fmt.Println("======================================== End Lab 3   ==========================================")
}
