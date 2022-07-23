package main

import (
	"fmt"
	"log"
	"time"
	"bytes"
	"net/http"
	"encoding/json"
	"encoding/base64"
)

const (
	Addr = "http://127.0.0.1:9000"
)

var Token string

func init () {
	strB := []byte("unencryptedDummyKey")
	encB := make([]byte, base64.StdEncoding.EncodedLen(len(strB))) 
	base64.StdEncoding.Encode(encB, strB)
	Token = string(encB)
}

type User struct {
	ID int `json:id`
	Name string `json:name`
	Password string `json:password`
}

type GenericResponse struct {
	Code int `json:code`
	Message string `json:message`
}

type SingleID struct {
	ID int `json:id`
}

func AddUser (user *User) {
	tmp, err := json.Marshal(user)
	if err != nil {
		log.Println(err)
	}
	b := bytes.NewReader(tmp)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("POST", Addr + "/add", b)
	if err != nil {
		log.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", Token)
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	
	// Read response code
	fmt.Println(res.StatusCode)
}

func GetUser (n int) {
	user := &User{}
	id := &SingleID{n}

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	bArray, err := json.Marshal(id)
	if err != nil {
		fmt.Println(err)
	}
	b := bytes.NewReader(bArray)

	req, err := http.NewRequest("GET", Addr + "/get", b)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", Token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Read response
	d := json.NewDecoder(res.Body)
	if res.StatusCode != 200 {
		fmt.Println(res.StatusCode)
		return
	}

	d.Decode(user)
	fmt.Printf("%d %s %s\n", user.ID, user.Name, user.Password)
}

func UpdateUser (id int, name, password string) {
	user := &User{id, name, password}
	byteArray, err := json.Marshal(user)
	if err != nil {
		fmt.Println(err)
	}

	reader := bytes.NewReader(byteArray)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("PUT", Addr + "/update", reader)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", Token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	
	// Read response code
	fmt.Println(res.StatusCode)
	return
}


func DeleteUser (n int) {
	id := &SingleID{n}
	byteArray, err := json.Marshal(id)
	if err != nil {
		fmt.Println(err)
	}
	reader := bytes.NewReader(byteArray)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("DELETE", Addr + "/delete", reader)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", Token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Read response code
	fmt.Println(res.StatusCode)
	return
}

func IncorrectRoute () {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", Addr + "/wrongpath", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Key", Token)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	// Read response code
	fmt.Println(res.StatusCode)
}


func main () {
	// Create users structs
	user1 := &User{1, "Bob", "12345678"}
	user2 := &User{2, "Ursella", "87654321"}
	user3 := &User{3, "Carl", "10101010"}
	
	// Post
	AddUser(user1)
	AddUser(user2)
	AddUser(user3)
	
	// Get
	GetUser(1)
	GetUser(2)
	GetUser(3)
	
	// Put
	UpdateUser(2, "Ursela", "")
	
	//Get
	GetUser(2)
	
	// Delete
	DeleteUser(3)
	DeleteUser(3)
	
	// Incorrect route
	IncorrectRoute()
}