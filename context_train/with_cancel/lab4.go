package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ExpriceServer(srvCTX context.Context, wg *sync.WaitGroup, chnl chan<- string, srvName string, duration time.Duration) {
	defer wg.Done()

	select {
	case <-srvCTX.Done():
		fmt.Printf("⏭️  [%s] Request cancelled (another server was faster) -> %v\n", srvName, srvCTX.Err())
		return
	case <-time.After(duration):
		select {
		case chnl <- fmt.Sprintf("[%s] (took %dms)", srvName, duration/time.Millisecond):
			fmt.Printf("\n🎉 [%s] responded first! Result sent.\n", srvName)
			return
		default:
			fmt.Printf("\n🎉 [%s] responded first! Result sent.\n", srvName)
		}
	}
}

func Lab4() {
	fmt.Println("======================================== Start Lab 4 ========================================")
	serverCTX, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	chnl := make(chan string, 1)
	defer close(chnl)

	servers := []string{"US-Server", "EU-Server", "Asia-Server"}
	durations := []time.Duration{200 * time.Millisecond, 600 * time.Millisecond, 800 * time.Millisecond}

	for index, srvName := range servers {
		wg.Add(1)
		fmt.Printf("🌐 Sending request to [%s]...\n", srvName)
		go ExpriceServer(serverCTX, wg, chnl, srvName, durations[index])
	}

	result := <-chnl
	cancel()
	wg.Wait()
	fmt.Printf("\n🏆 Fastest result: Success from %s\n\n", result)

	fmt.Println("✅ All competitive requests finished cleanly. No goroutine leaks")

	fmt.Println("======================================== End Lab 4   ==========================================")

}
