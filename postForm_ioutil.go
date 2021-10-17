package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Post struct {
	Userid string `json:"userId"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func main() {

	data := url.Values{}
	data.Add("title", "sunt aut facere repellat provident occaecati excepturi optio reprehenderit")
	data.Add("body", "HTTP test server accepting GET/POST requests")
	data.Add("userId", "999")
	data.Add("id", "9")

	resp, err := http.PostForm("https://jsonplaceholder.typicode.com/posts", data)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))
	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Posted title:", post.Title)

}
