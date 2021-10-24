package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Post struct {
	UserId int    `json:"userId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}

var counter int
var mutex = &sync.Mutex{}

func process(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	if id, err := strconv.Atoi(vars.Get("Id")); err == nil {
		log.Printf("Response a single post by Id = %d\n", id)
		idHandler(id, w, r)
		return
	}
	keys, ok := r.URL.Query()["Name"]
	if ok && len(keys[0]) >= 1 {
		log.Println("Response multiple posts by Name = " + string(keys[0]))
		paraNameHandler(w, r)
		return
	}
	log.Println("Response all users")
	usersHandler(w, r)
}
func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method is not POST.", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "POST request successful\n")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	//http://localhost:8082/users
	posts := readJsonFile()
	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func idHandler(id int, w http.ResponseWriter, r *http.Request) {
	//http://localhost:8082/users?Id=1
	posts := readJsonFile()
	if (id >= len(posts)) || (id < 0) {
		http.Error(w, "Requested index is out of range!", http.StatusInternalServerError)
		return
	}
	js, err := json.Marshal(posts[id])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func paraNameHandler(w http.ResponseWriter, r *http.Request) {
	//http://localhost:8082/users?Name=
	posts := readJsonFile()
	js, err := json.Marshal(posts[1:3])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func readJsonFile() []Post {
	filename, err := os.Open("users.json")
	if err != nil {
		log.Fatal(err)
	}
	defer filename.Close()
	data, err := ioutil.ReadAll(filename)
	if err != nil {
		log.Fatal(err)
	}
	var result []Post
	jsonErr := json.Unmarshal(data, &result)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	return result
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, "ping http://localhost:8082 count:%s", strconv.Itoa(counter))
	mutex.Unlock()
}
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!\n", r.URL.Path[1:])
	})

	http.HandleFunc("/users", process)

	http.HandleFunc("/form", postHandler)

	// response ping for fun
	http.HandleFunc("/ping", pingHandler)

	fmt.Printf("Starting server at port 8082...\n")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
