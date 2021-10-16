package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// create a get request and issue it by Get function
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	// check if error
	if err != nil {
		log.Fatal(err)
	}
	// Client must always close the response body when finished
	defer resp.Body.Close()
	// copy response body to os standard out, return number of bytes copied and the first error occured while copying
	copiedBytes, err := io.Copy(os.Stdout, resp.Body)
	// check if error encoutered
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("number of bytes copied: ", copiedBytes)
}
