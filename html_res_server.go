package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var counter int
var mutex = &sync.Mutex{}

/*func httpPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, r.URL.Path[1:])
}*/

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "POST request successful\n\n")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	fmt.Fprintf(w, "First name = %s\n", fname)
	fmt.Fprintf(w, "Last name = %s\n", lname)
}

func plainHtmlHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Response HTML plain text</h1>"+
		"<form action=\"/form\" method=\"POST\">"+
		"<label for=\"fname\">First name:</label><br>"+
		"<input type=\"text\" id=\"fname\" name=\"fname\" value=\"John\"><br>"+
		"<label for=\"lname\">Last name:</label><br>"+
		"<input type=\"text\" id=\"lname\" name=\"lname\" value=\"Doe\"><br><br>"+
		"<input type=\"submit\" value=\"Submit\">"+
		"</form>")
}
func pingHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, "ping http://localhost:8081 count:%s", strconv.Itoa(counter))
	mutex.Unlock()
}
func main() {
	// response some text
	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from http://localhost:8081")
	})
	// response plain html
	http.HandleFunc("/plainHtml", plainHtmlHandler)
	// response a static file
	http.Handle("/", http.FileServer(http.Dir("./static")))
	// handle POST request
	http.HandleFunc("/form", formHandler)
	// response ping for fun
	http.HandleFunc("/ping", pingHandler)
	fmt.Printf("Starting server at port 8081...\n")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
