package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	//Encode the data
	postBody, err := json.Marshal(map[string]string{
		"userId": "100",
		"id":     "101",
		"title":  "ping ping ping",
		"body":   "Happiness is an attitude. We either make ourselves miserable, or happy and strong. The amount of work is the same",
	})
	if err != nil {
		log.Fatal(err)
	}
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", responseBody)

	if err != nil {
		log.Fatalf("An error occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf(string(body))
}
