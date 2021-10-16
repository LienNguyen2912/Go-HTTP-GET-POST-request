package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// create a Get request to a webpage then issue it by Get function
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	// check for the error
	if err != nil {
		log.Fatal(err)
	}
	// the client must close the response body when finished
	defer resp.Body.Close()
	// read all the response body by ReadAll
	body, err := ioutil.ReadAll(resp.Body)
	// check for error
	if err != nil {
		log.Fatal(err)
	}
	// print the received content
	fmt.Println(string(body))
}
