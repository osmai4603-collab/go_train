package main

import (
	"context"
	"fmt"
	"sync"
)

func SearchFailure(ctx context.Context, wg *sync.WaitGroup) {

}

func Lab5() {
	fmt.Println("======================================== Start Lab 5 ========================================")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	SearchFailure(ctx, wg)

	fmt.Println("======================================== End   Lab 5 ========================================")
}
