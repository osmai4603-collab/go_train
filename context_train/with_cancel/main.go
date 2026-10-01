package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ==============================================================================
// 1. Example 1: Worker Pool Graceful Shutdown
// ==============================================================================
// In this scenario, we spawn multiple goroutines to process continuous tasks.
// When the supervisor (main) decides to stop them, it calls cancel().
// All workers receive the signal via <-ctx.Done() and exit cleanly (no goroutine leak).
func runWorkerPoolExample() {
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("📌 Scenario 1: Worker Pool Graceful Shutdown using context.WithCancel")
	fmt.Println("------------------------------------------------------------------")

	// 1. Create a cancellable context derived from context.Background()
	// ctx: context passed down to workers to listen for cancellation
	// cancel: function that closes ctx.Done() channel when invoked
	ctx, cancel := context.WithCancel(context.Background())

	// Best Practice: Always defer cancel() to release resources and prevent leaks
	defer cancel()

	var wg sync.WaitGroup

	// Launch 3 workers
	workerCount := 3
	for i := 1; i <= workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("  [Worker %d] 🟢 Started working and monitoring...\n", workerID)

			select {
			case <-ctx.Done():
				// When cancel() is called, ctx.Done() is closed immediately.
				// ctx.Err() indicates why it was cancelled (context.Canceled).
				fmt.Printf("  [Worker %d] 🛑 Received cancel signal (Reason: %v). Shutting down...\n",
					workerID, ctx.Err())
				return

			default:
				// Simulating active work
				fmt.Printf("  [Worker %d] ⚙️  Processing a subtask...\n", workerID)
				//time.Sleep(4 * time.Second)
			}
		}(i)
	}

	// Allow workers to do some work for 1 second
	time.Sleep(3 * time.Second)

	fmt.Println("\n📢 [Supervisor]: Trigger event occurred, calling cancel() to stop all workers!")
	fmt.Println()
	// Calling cancel() closes ctx.Done(), broadcasting to all listening goroutines simultaneously
	cancel()

	// Wait for all workers to cleanly shut down
	wg.Wait()
	fmt.Println("✅ [Supervisor]: All workers stopped gracefully without any goroutine leaks.")
	fmt.Println()
}

// ==============================================================================
// 2. Example 2: Competitive Search / Race & Early Exit
// ==============================================================================
// In this scenario, we query multiple replica servers in parallel.
// As soon as the fastest server responds with the needed data,
// we immediately invoke cancel() to terminate the other requests,
// saving network bandwidth and CPU cycles.
func runFastestResponseExample() {
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("📌 Scenario 2: Competitive Search & Early Cancellation (Race & Early Exit)")
	fmt.Println("------------------------------------------------------------------")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	resultCh := make(chan string, 1)

	// Simulating 3 replica servers with different response times
	servers := []struct {
		name     string
		duration time.Duration
	}{
		{name: "US-Server", duration: 800 * time.Millisecond},
		{name: "EU-Server", duration: 300 * time.Millisecond}, // Fastest
		{name: "Asia-Server", duration: 600 * time.Millisecond},
	}

	for _, s := range servers {
		wg.Add(1)
		go func(serverName string, delay time.Duration) {
			defer wg.Done()
			fmt.Printf("  🌐 Sending request to [%s]...\n", serverName)

			select {
			case <-time.After(delay):
				// This server finished first
				select {
				case resultCh <- fmt.Sprintf("Success from [%s] (took %v)", serverName, delay):
					fmt.Printf("  🎉 [%s] responded first! Cancelling other requests...\n", serverName)
					cancel() // Cancel remaining server requests immediately!
				default:
				}

			case <-ctx.Done():
				// Cancelled because another server already won the race
				fmt.Printf("  ⏭️  [%s] Request cancelled (another server was faster).\n", serverName)
			}
		}(s.name, s.duration)
	}

	// Read the winning result
	fastestResult := <-resultCh
	fmt.Printf("\n🏆 Received result in main(): %s\n\n", fastestResult)

	// Wait for all goroutines to finish clean exit
	wg.Wait()
	fmt.Println("✅ All competitive tasks finished cleanly.")
	fmt.Println()
}

// ==============================================================================
// 3. Example 3: Fail-Fast on Error (إلغاء كافة المهام عند فشل إحداها)
// ==============================================================================
// In this scenario, multiple subtasks run in parallel to process an order.
// If any subtask encounters an error, it calls cancel() immediately.
// All remaining running subtasks abort via <-ctx.Done(), saving resources.
func runErrorCancellationExample() {
	fmt.Println("------------------------------------------------------------------")
	fmt.Println("📌 Scenario 3: Fail-Fast on Error (إلغاء المهام فور حدوث خطأ في إحداها)")
	fmt.Println("------------------------------------------------------------------")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var errOnce sync.Once
	var firstErr error

	tasks := []struct {
		id       int
		name     string
		duration time.Duration
		err      error
	}{
		{id: 1, name: "Validate User Account", duration: 200 * time.Millisecond, err: nil},
		{id: 2, name: "Process Payment", duration: 500 * time.Millisecond, err: fmt.Errorf("payment declined (insufficient funds)")},
		{id: 3, name: "Reserve Inventory", duration: 1200 * time.Millisecond, err: nil},
		{id: 4, name: "Generate Shipping Label", duration: 1500 * time.Millisecond, err: nil},
	}

	for _, task := range tasks {
		wg.Add(1)
		go func(id int, name string, delay time.Duration, taskErr error) {
			defer wg.Done()
			fmt.Printf("  [Task %d: %s] ⏳ Processing started...\n", id, name)

			select {
			case <-time.After(delay):
				if taskErr != nil {
					// Task failed: capture the error and cancel all other tasks
					errOnce.Do(func() {
						firstErr = taskErr
						fmt.Printf("  ❌ [Task %d: %s] Failed: %v\n", id, name, taskErr)
						fmt.Println("  📢 Calling cancel() to abort remaining running tasks!")
						cancel()
					})
					return
				}
				fmt.Printf("  ✅ [Task %d: %s] Completed successfully.\n", id, name)

			case <-ctx.Done():
				// Another task failed and called cancel() while this task was still executing
				fmt.Printf("  🛑 [Task %d: %s] Aborted early due to context cancellation: %v\n", id, name, ctx.Err())
			}
		}(task.id, task.name, task.duration, task.err)
	}

	// Wait for all goroutines to finish or cleanly abort
	wg.Wait()

	fmt.Println()
	if firstErr != nil {
		fmt.Printf("⚠️ Transaction Aborted: Failure occurred -> %v\n\n", firstErr)
	} else {
		fmt.Println("🎉 Transaction Completed: All subtasks succeeded.")
		fmt.Println()
	}
}

// func main() {
// 	// fmt.Println("==================================================================")
// 	// fmt.Println("🚀 Comprehensive Guide to context.WithCancel in Go")
// 	// fmt.Println("==================================================================")
// 	// fmt.Println()

// 	// // Run first scenario
// 	// runWorkerPoolExample()

// 	// // Short pause between examples
// 	// time.Sleep(500 * time.Millisecond)

// 	// Run second scenario
// 	runFastestResponseExample()

// 	// Short pause between examples
// 	time.Sleep(500 * time.Millisecond)

// 	// // Run third scenario (Fail-Fast on Error)
// 	// runErrorCancellationExample()

// 	// fmt.Println("==================================================================")
// 	// fmt.Println("🎉 All examples completed successfully!")
// 	// fmt.Println("==================================================================")
// }

func main() {
	// Lab1()
	// Lab2()
	// Lab3()
	Lab4()
}
