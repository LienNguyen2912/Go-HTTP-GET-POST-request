package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func main() {

	urls := []string{
		"http://localhost:8081/ping",
		"http://localhost:8082/ping",
	}

	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		for _, url := range urls {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				content := doGetRequest(url)
				fmt.Println(content)
			}(url)
		}
	}
	wg.Wait()
}

func doGetRequest(url string) (content string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}
