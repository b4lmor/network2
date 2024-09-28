package utils

import (
	"fmt"
	"sync"
)

type SafePrinter struct {
	mu sync.Mutex
}

func (sp *SafePrinter) Print(message string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	fmt.Println(message)
}
