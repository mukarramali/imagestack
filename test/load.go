package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var IMAGE_STACK_URL = "https://imagestack-latest.sliplane.app"

func downloadFile(url string, wg *sync.WaitGroup) {
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching:" + url)
	}
	resp.Body.Close()
	wg.Done()
}

func generateRandomUrl() string {
	quality := rand.Intn(100) + 1
	picSumId := rand.Intn(100) + 1
	return fmt.Sprintf("%s/?quality=%d&url=https://picsum.photos/id/%d/5000", IMAGE_STACK_URL, quality, picSumId)
}

func main() {
	var wg sync.WaitGroup
	tests := 10

	if len(os.Args) > 1 {
		tests, _ = strconv.Atoi(os.Args[1])
	}

	for i := 1; i <= tests; i++ {
		wg.Add(1)
		go downloadFile(generateRandomUrl(), &wg)
	}

	wg.Wait()
}
