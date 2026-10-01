package main

import (
	"context"
	"fmt"
	"time"
)

func withTimeout() {
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 3*time.Second)
	defer cancel()

	deadline, ok := ctx.Deadline()

	if ok {
		fmt.Println("Deadline: True")
	}

	fmt.Printf("time rest: %w\n", deadline)

	select {
	case <-parent.Done():
		fmt.Println("context timeout")
	default:
		fmt.Println("Stil waiting...")
	}
}

func WithCancel() {
	parent := context.Background()

	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Cancellation")
				return
			default:
				fmt.Println("working")
				time.Sleep(5 * time.Second)
			}
		}
	}()
	time.Sleep(15 * time.Second)
}

func main() {
	withTimeout()
}
