# Go HTTP Client/Server
## GET
Make a http GET request to https://jsonplaceholder.typicode.com/posts/1 and print out the result
```sh
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
```
The values returned from _http.Get_ typically called _resp_ and _err_. If there was an error happenned, err will be non nil. _log.Fatal(err)_ is equivalent with 

> log.Print(err)</br>
> os.Exit(1)

means it prints the message then terminates the program.
_resp_ contains not only the data which we are interested in but also incoherent data like the header and properties of the request made. So that we have to access the _Body_ property and read it by _ioutil.ReadMe_ invocation. 

### ✨ Remark: It is important to close the response Body to prevent memory leaks.<br>
The program above will print</br>
![get1](https://user-images.githubusercontent.com/73010204/137583729-3ef9117f-7531-40fb-93b3-dfca5bd0aa49.PNG)</br>
Let's try a non existed http address to get error.
```sh
resp, err := http.Get("https://blablabla.com")
```
Here it is</br>
![getErr](https://user-images.githubusercontent.com/73010204/137584208-7e9a9470-5f03-4581-8840-44e3fccb1830.PNG)</br>
Instead of _ioutil.ReadAll_ we can use the _io.Copy_ function.</br>
The next program uses _io.Copy_ to copy the response body to os standard out. It returns number of bytes copied and the first error occured while copying, if any.
```sh
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
```
and the output is same.</br>
![getIoCopy](https://user-images.githubusercontent.com/73010204/137584618-2bed6d24-7194-4347-9ea9-970a955b037c.PNG)</br>

## GET request with query parameters
Say we want to issue https://jsonplaceholder.typicode.com/comments?postId=2 request. Query parameter is "postId=2".
```sh
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
```
The output has 5 records as below</br>
![get_para](https://user-images.githubusercontent.com/73010204/137586056-e801baed-b998-4b4e-a3f6-95fab2428570.PNG)</br>

## Custom GET request
In the previous pgrogram, we have used the _http.Get_ which is simple and quickly make GET requests. But what if we want to custom something like add a timeout? _http.Client_ comes to help.</br>
```sh
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("GET", "https://jsonplaceholder.typicode.com/users", nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(body))
}
```
We created our own _http.Client_, specify Timeout property, then call _.Do(request)_. </br>
How about custom request having query parameters? Access _URL.Query_ property of request then call _Add_ method, see the example belows to issue https://jsonplaceholder.typicode.com/comments?postId=2 Get request.</br>
_We also add  **_json.Unmarshal_** invocation to get the first record only, for fun!_
```sh
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
```
Here the output is</br>
![customeAddMarshal](https://user-images.githubusercontent.com/73010204/137586594-7119a763-0bba-48c3-91dd-8c62b2390d2f.PNG)</br>

## Go POST 
The _http.Post_ function takes three parameters.
- The URL address of the server
- The content type of the body as a string, for example, _application/json_
- The request body whose type is _io.Reader_. To create a request body from json format, we need
	- Encode JSON data by _Marshall_ invocation. We will get a []bytes and error non nil (if any)
	- Convert the encoded JSON data to a type implemented by the _io.Reader_ interface by using _bytes.NewBuffer_ function

```sh
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
```
The program will print</br>
![post1](https://user-images.githubusercontent.com/73010204/137624048-554e973f-9d87-4f1a-8998-cf94d6f6218b.PNG)</br>
## Go POST FORM
We can use the _http.PostForm_ to submit a form easily. 
```sh
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
```
The program will print</br></br>
![postForm](https://user-images.githubusercontent.com/73010204/137624586-663716a2-2868-46e1-9323-db05ed2ae693.PNG)</br></br>
We could use _ioutil.ReadAll_ like the examples above to first read the data into memory and then call _json.Unmarshall_. 
```sh
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
```
We 'll get</br>
![postForm_ioutil](https://user-images.githubusercontent.com/73010204/137626123-251fb547-0385-468b-a488-f7dfcd54d8c9.PNG)</br>
That's it!
## Create a simple http server
## Rendering Plain Text - Serving a File - Respond to a POST request with form submission
```sh
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
```
In the example above, we’ve created our own http://localhost:8081 server. Execute it and hopefully see it works</br>
- Rendering some text _http://localhost:8081/hi_</br>
![hi](https://user-images.githubusercontent.com/73010204/138557099-8e9e5327-be8c-4195-8dbf-d469005501bc.PNG)</br>
- Rendering a plain html content _http://localhost:8081/plainHtml_</br>
![plainTHML](https://user-images.githubusercontent.com/73010204/138557101-f8afa422-f70c-43fa-8381-833d9713a76e.PNG)</br>
- Handling a Post request and rendering some text _http://localhost:8081/form_</br>
![postRes](https://user-images.githubusercontent.com/73010204/138557104-b12cec75-98d6-4a6c-a046-44a1521cbaa7.PNG)</br>
- Rendering a static html file<br>
![html](https://user-images.githubusercontent.com/73010204/138557115-022668d7-753e-46cd-8247-5c2b01c629dd.PNG)<br>
Here is the folder structure I used</br>
![staticF](https://user-images.githubusercontent.com/73010204/138557706-751fa155-5533-46dd-8332-461c37ea9d47.png)</br>
- Finally, adding a mutex implementation for fun! Whenever user accesses _http://localhost:8081/ping_ the server responses a number increased by 1 continuously</br>
![ping](https://user-images.githubusercontent.com/73010204/138557120-eaf8712d-040e-470f-b785-0c6618e7f68a.PNG)</br>

**What to notice:**
- http.ListenAndServer() starts the server and listens on the TCP network address :8081. This function blocks blocks the current goroutine.
- ListenAndServe specifies the port to listen on as the first argument and an http.Handler as its second argument. If the handler is nil, DefaultServeMux is used.
>While DefaultServeMux is okay for toy programs, you should avoid it in your production code. This is because it is accessible to any package that your application imports, including third-party packages. Therefore, it could potentially be exploited to expose a malicious handler to the web if a package becomes compromised.
So, what's the alternative? A locally scoped http.ServeMux!
Refer https://www.honeybadger.io/blog/go-web-services/
- If there are a lot of static HTML files created by hand, _html/template_ package would be prefered.
## Rendering JSON
Firstly, we need to prepare json data. There are many ways to do it. We can create a variable and assign json data to it as below
```sh
type Post struct {
	PostId int    `json:"postId"`
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}
func createPosts() []Post {
	PostsSlice := make([]Post, 0, 2)
	PostsSlice = append(PostsSlice, Post{
		PostId: 0,
		Id:     0,
		Name:   "eaque et deleniti atque tenetur ut quo ut",
		Email:  "Hautdfsfsfsden@althea.biz",
		Body:   "voluptate iusto quis nobis reprehenderit ipsum amet nulla\nquia quas dolores velit et non\naut quia necessitatibus\nnostrum quaerat nulla et accusamus nisi facilis",
	})
	PostsSlice = append(PostsSlice, Post{
		PostId: 1,
		Id:     1,
		Name:   "odio adipisci rerum aut animi",
		Email:  "Nikita@garfield.biz",
		Body:   "uia molestiae reprehenderit quasi aspernatur\naut expedita occaecati aliquam eveniet laudantium\nomnis quibusdam delectus saepe quia accusamus maiores nam est\ncum et",
	})
	return PostsSlice
}
func usersHandler(w http.ResponseWriter, r *http.Request) {
	posts := createPosts()
	js, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func main() {
	http.HandleFunc("/users", usersHandler)

	fmt.Printf("Starting server at port 8082...\n")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
```
Or we can read JSON data from a file. </br>
This is users.json data. Refer https://jsonplaceholder.typicode.com/posts
```
[
    {
    "userId": 0,
    "id": 0,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
  },
  {
    "userId": 1,
    "id": 1,
    "title": "qui est esse",
    "body": "est rerum tempore vitae\nsequi sint nihil reprehenderit dolor beatae ea dolores neque\nfugiat blanditiis voluptate porro vel nihil molestiae ut reiciendis\nqui aperiam non debitis possimus qui neque nisi nulla"
  }
]
```
Next is our simple server 
```sh
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
	posts := readJsonFile()
	js, err := json.Marshal(posts[1:3])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func pingHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	counter++
	fmt.Fprintf(w, "ping http://localhost:8082 count:%s", strconv.Itoa(counter))
	mutex.Unlock()
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
```
Let's execute it and hopefully see it works
- Response to GET request without parameter _http://localhost:8082/users_</br>
![users](https://user-images.githubusercontent.com/73010204/138593371-75af3fe4-9d1f-4db4-9832-c6549e1eaf1b.PNG)</br>
- Response to GET request with parameter _http://localhost:8082/users?Id=1_</br>
![userId](https://user-images.githubusercontent.com/73010204/138593458-ad684c08-9c38-476f-9ce8-fffd8836cb63.PNG)</br>
- If there is POST request, _http://localhost:8082/form , just rendering a "POST request successful" message. 
- Finally, adding a mutex implementation for fun! Whenever user accesses _http://localhost:8082/ping_ the server responses a number increased by 1 continuously</br>
![pingping](https://user-images.githubusercontent.com/73010204/138593548-8b96ae80-7448-42ff-8515-89d8ddfddfaf.PNG)</br>
## Async GET requests
To the multiple local servers above, say we create multiple asynchronous HTTP requests to them.
```sh
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
```
The output may vary each time executing.</br>
![asynGet](https://user-images.githubusercontent.com/73010204/138594181-baf7c251-30dc-437e-8ee3-de9edc44db5e.PNG)</br>
## Summary
We have implemented a simple HTML client/server that
- Issue GET requests
- Issue GET requests having parameter
- Issue POST requests
- Customize HTML requests
- Rendering a static file
- Rendering plain HTML text
- Rendering Json data
- Retrieve paramaters from a request
- Reponse to a POST request and parser POST form
- Async GET request.
