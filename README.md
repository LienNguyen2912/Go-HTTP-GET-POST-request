# Go HTTP GET/POST request
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


