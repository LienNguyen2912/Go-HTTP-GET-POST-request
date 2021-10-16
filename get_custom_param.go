package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Posts []struct {
	PostId int    `json:"postId"`
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

func main() {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/comments", nil)
	if err != nil {
		log.Fatalln(err)
	}
	// create https://jsonplaceholder.typicode.com/comments?postId=2 request
	q := request.URL.Query()
	q.Add("postId", "2")
	request.URL.RawQuery = q.Encode()
	fmt.Println(request.URL.String())

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// log.Println(string(body)) // not print all the body received
	// Unmarshal result
	posts := Posts{}
	err = json.Unmarshal(body, &posts)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Name of 1st record: %s", posts[0].Name)
}
