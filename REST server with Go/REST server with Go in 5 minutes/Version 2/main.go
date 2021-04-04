package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

var userlist []User

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

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
func postUser(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint Hit: postUser")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var userData User
	json.Unmarshal(reqBody, &userData)
	userlist = append(userlist, userData)
	JSON(w, http.StatusCreated, "Created")
}
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
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/user/get/", getUser)
	mux.HandleFunc("/user/post/", postUser)
	mux.HandleFunc("/user/delete/", deleteUser)
	mux.HandleFunc("/user/patch/", patchUser)
	log.Fatal(http.ListenAndServe(":8090", mux))
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}
}
func getUid(pathSuffix string, r *http.Request) int {
	var slug string
	if strings.HasPrefix(r.URL.Path, pathSuffix) {
		slug = r.URL.Path[len(pathSuffix):]
	}
	uid, _ := strconv.Atoi(slug)
	return uid

}
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
