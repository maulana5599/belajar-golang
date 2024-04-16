package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

var (
	data  map[string]string
	mutex sync.Mutex
)

func writeToMap(key, value string) {
	data[key] = value
	fmt.Println(value)
}

func main() {
	data = make(map[string]string)

	runtime.GOMAXPROCS(16)

	for i := 0; i < 10; i++ {
		stringNumber := fmt.Sprintf("%d", i)
		go func() {
			writeToMap(stringNumber, stringNumber)
		}()
	}

	// Menunggu sedikit agar goroutine memiliki kesempatan untuk menjalankan
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Hello from main!")
	fmt.Println(data)
}
