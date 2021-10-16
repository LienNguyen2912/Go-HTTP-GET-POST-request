package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func main() {
	// create https://jsonplaceholder.typicode.com/comments?postId=2 request
	params := url.Values{}
	params.Add("postId", "2")
	resp, err := http.Get("https://jsonplaceholder.typicode.com/comments?" + params.Encode())

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(body))
}
