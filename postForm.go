package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {

	data := url.Values{
		"userId": {"1"},
		"id":     {"3"},
		"title":  {"sunt aut facere repellat provident occaecati excepturi optio reprehenderit"},
		"body":   {"HTTP test server accepting GET/POST requests"},
	}

	resp, err := http.PostForm("https://httpbin.org/post", data)

	if err != nil {
		log.Fatal(err)
	}

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res["form"])
}
