👋🏻Hey, If the title of this article excites you, then my friend, you're at right place.

Let's quickly go through what we're going to cover in this post.

1. First will be going to build a teeny weeny REST server and understand its working. After that, we'll incrementally improve its working and know its weakness.
2. We build a CRUD server where we can `CREATE`, `READ`, `UPDATE` and `DELETE` User details.

To keep everything simple. We'll focus on the basic concepts and won’t interact with any package and build server using the standard libraries and understand their limit and strength.

## Prerequisites

You'll need [Go version 1.16+](https://golang.org/dl/) installed on your development machine.

## Start 👷🏻‍♂️

Now let's dive into the code.

### Version v1.0

```go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {

	log.Println("Endpoint Hit: homePage")

	fmt.Fprintf(w, "Welcome to the HomePage!")
}
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8090", mux))
}

```

If you run this code on your machine. This will start simple API server at port `8090`. In your browser navigate to `http://localhost:8090/`, you'll going to see "Welcome to the HomePage!".

![Output](https://github.com/KushagraMehta/Blog/blob/master/REST%20server%20with%20Go/REST%20server%20with%20Go%20in%205%20minutes/Version_1-Output.png)

## THE API IS CREATED 🎇

![Ta-Da](https://media.giphy.com/media/Y3RqktJpqtVl9bs4Ee/giphy.gif)

Ta-Da, You build your first API server with Go. Yes, It was that easy.

Ok, Now understand the code

```go
func http.NewServeMux() *http.ServeMux
```

> `http.NewServeMux` is the default serve mux in Go from the **net/http** package. The method will actually create a multiplexer where it registers patterns with their corresponding function. We can also skip this step, as by default **net/http** package create `DefaultServeMux` which is the same.

```go
func (*http.ServeMux).HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
```

> The `HandleFunc` method takes two arguments, the first _string pattern_, and the second _handler function_ to which pattern to register with. Ok, but what the heck is a handler. Handler is an interface that has a method called _ServeHttp_, responds to an HTTP request, and has all the logic related to that End-point. We'll understand more about it in down the post.

```go
func ListenAndServe(addr string, handler Handler) error
```

> Then `ListenAndServe` method takes two arguments, the first is to address to start the server on and handler in the case of using **DefaultServeMux** we can pass nill as an argument.

### Version v2.0

Now we're going to build REST server with all the functionality which will allow us to `CREATE`, `READ`, `UPDATE` and `DELETE` users.

#### User Structure

Before starting with anything first build the user structure and list in which our user will going to store in.

```go
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// let's declare a global Articles array
// that we can then populate in our main function
// to simulate a database
var userlist []User
```

You can see that we have use tags on struct field declarations to customize the encoded JSON key names. These tags are going to be used by `"encoding/json"` package.

#### Create Helper function

**JSON** function is going to encode any data to JSON format so that we can send the data to the client.

```go
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
```

**getUid** will extract _user-id(uid)_ from **URL** by slicing the **URL**. After that, it converts the _uid_ to int from the string.

```go
func getUid(pathSuffix string, r *http.Request) int {
	var slug string
	if strings.HasPrefix(r.URL.Path, pathSuffix) {
		slug = r.URL.Path[len(pathSuffix):]
	}
	uid, _ := strconv.Atoi(slug)
	return uid
}
```

**Error** is just a wrapper around **JSON** function to add an additional error flag to the JSON data

```go
func Error(w http.ResponseWriter, statusCode int, err error) {
	if err != nil {
		JSON(w, statusCode, struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	JSON(w, http.StatusBadRequest, nil)
}
```

#### Building our Router

Now we'll add our routes to `main.go` file.

```go
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/user/get/", getUser)
	mux.HandleFunc("/user/post/", postUser)
	mux.HandleFunc("/user/delete/", deleteUser)
	mux.HandleFunc("/user/patch/", patchUser)

	log.Fatal(http.ListenAndServe(":8090", mux))
}
```

#### Create

We will need to create a new function that will do the job of creating new users.

Let’s start off by creating a `postUser()` function within our `main.go` file.

```go
func postUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: postUser")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userData User
	json.Unmarshal(reqBody, &userData)
	userlist = append(userlist, userData)
	JSON(w, http.StatusCreated, "Created")
}
```

Now understand what we have written in the function.

1. First of all, we've to extract the bytes array from the request body, for that we use `ioutil.ReadAll(r.Body)` which reads from **r** until an error or _EOF_ and returns the data it read.
2. `json.Unmarshal(reqBody, &userData)` parses the JSON-encoded data and stores the result in the .userData.
3. After that we append the user into userlist.
4. Send Created respone to the client

#### Read

Now we have created the creation function so now we want to fetch the newly created user from the server.

For that we'll create a `getUser()` function within our `main.go` file.

```go
func getUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: getUser")
	uid := getUid("/user/get/", r)
	for _, user := range userlist {
		if user.ID == uid {
			JSON(w, http.StatusOK, user)
			return
		}
	}
	Error(w, http.StatusNotFound, errors.New("User Not Found"))
}
```

Now understand what we have written in the function.

1. As client will send the user id in _URL_ so we parse the _URL_ and extract uid.
2. After that we'll loop over the userlist array to find the user
   - If we find one then we send OK status with user data.
   - If the user is not present then send an Error response that the user is not found.

#### Update

After creating and reading now it's time for updating. create a `patchUser()` function within our `main.go` file.

```go
func patchUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: patchUser")
	uid := getUid("/user/patch/", r)
	reqBody, _ := ioutil.ReadAll(r.Body)
	var updatedUserData User
	json.Unmarshal(reqBody, &updatedUserData)
	for index, user := range userlist {
		if user.ID == uid {
			userlist[index] = updatedUserData
			JSON(w, http.StatusOK, "Patched")
			return
		}
	}
	Error(w, http.StatusNotFound, errors.New("User Not Found"))
}
```

#### Delete

There may be times where you need to delete the data being exposed by your REST API. In order to do this, you need to expose a DELETE endpoint within your API that will take in an identifier and delete whatever is associated with that identifier.

Add a new function to your main.go file which we will call `deleteUser()`:

```go
func deleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: deleteUser")
	uid := getUid("/user/delete/", r)
	for index, user := range userlist {
		if user.ID == uid {
			userlist = append(userlist[:index], userlist[index+1:]...)
			JSON(w, http.StatusOK, "Deleted")
			return
		}
	}
	Error(w, http.StatusNotFound, errors.New("User Not Found"))
}
```

> Note - To keep this simple, we are updating a global variable. However, we aren’t doing any checks to ensure that our code is free of race conditions. In order to make this code thread-safe, I recommend checking out [Go Mutexes](https://gobyexample.com/mutexes)

if you're reading this, that means now you can build any CRUD application.

![](https://media.giphy.com/media/3oz8xDLuiN1GcDA3xC/giphy.gif)

## Conclusion

This example represents a very simple RESTful API written using Go. In a real project, we’d typically tie this up with a database so that we were returning real values. For the next step, I would suggest you should read [Gorilla/Mux](https://github.com/gorilla/mux) package. As it will remove a lot of boilerplate code and written only with the standard line.

> Source Code - The full source code for this tutorial can be found here: [KushagraMehta/Blog/REST server with Go](https://github.com/KushagraMehta/Blog/tree/master/REST%20server%20with%20Go/REST%20server%20with%20Go%20in%205%20minutes)

In the next post, we'll add Postgresql plus Dockerize the whole application, furthermore host it on Heroku/AWS.

I hope you find this blog useful. Please share your thought in the comments.
![](https://media.giphy.com/media/Ely3cY20BUEu8rWtDH/giphy.gif)
