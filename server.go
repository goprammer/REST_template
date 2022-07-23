// Package handlers Example RESTful Server with Swagger notation.
//
// Standard example of turning go code into swagger.
//
// Schemes: http
// Host: localhost
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"log"
	"time"
	"net/http"
	"encoding/json"
	"encoding/base64"

	"github.com/gorilla/mux"
	"github.com/go-openapi/runtime/middleware"
)

// @termsOfService http://swagger.io/terms/

var allUsers = map[int]*SingleUser{}

var Token string

func init () {
	strB := []byte("unencryptedDummyKey")
	encB := make([]byte, base64.StdEncoding.EncodedLen(len(strB))) 
	base64.StdEncoding.Encode(encB, strB)
	Token = string(encB)
}

// SingleUser
// swagger:model
type SingleUser struct {
	// The ID for the user
	// in: body
	//
	// required: true
	// min: 1
	ID int `json:"id"`

	// The Name of the user
	// in: body
	//
	// required: true
	Name string `json:"name"`

	// The Password of the user
	//
	// required: true
	Password string `json:"password"`
}

// SingleID
// swagger:model
type SingleID struct {
	// The ID for the user
	// in: body
	//
	// required: true
	// min: 1
	ID int `json:"id"`
}


// Success: Your requested content has been sent in response body. HTTP status code 200 in header.
// swagger:response OK
type OK struct {
	// in:body
	Code int
}


// Success: Created user. HTTP status code 201 in header.
// swagger:response Created
type Created struct {
	// in:body
	Code int
}

// Success: Completed request, but no content to return. HTTP status code 204 in header.
// swagger:response NoContent
type NoContent struct {
	// in:body
	Code int
}

// Error: Couldn't read client's json. HTTP status code 400 in header.
// swagger:response BadRequest
type BadRequest struct {
	// in:body
	Code int
}

// Error: API creditials have been rejected. HTTP status code 401 in header.
// swagger:response Unauthorized
type Unauthorized struct {
	// in:body
	Code int
}

// Error: Couldn't find user. HTTP status code 404 in header.
// swagger:response NotFound
type NotFound struct {
	// in:body
	Code int
}

// Error: User already exists. HTTP status code 409 in header.
// swagger:response Conflict
type Conflict struct {
	// in:body
	Code int
}

// swagger:response UserResponse
type UserResponse struct {
	// Return information of a user.
	// in:body
	Body SingleUser
}


// swagger:parameters UserParam AddUserParam
type UserParam struct {
	// Requested information for a single user.
	// in:body
	Body SingleUser
}

// swagger:parameters getUser deleteUser
type IDParam struct {
	// Client provides ID for user it wants info for.
	// in:body
	// required:true
	Body SingleID
}


// GetSingleUser returns a single user matched to id provided by client.
func GetSingleUser (w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /get GetSingleUser getUser
	// Matches id to find a single user.
	//
	// Responses:
	// 200: UserResponse
	// 400: BadRequest
	// 401: Unauthorized
	// 404: NotFound
	auth:= r.Header.Get("X-Auth-Key")
	if auth != Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := &SingleID{}
	err := json.NewDecoder(r.Body).Decode(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	user, ok := allUsers[id.ID]
	if ok {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
		return
	}
	
	w.WriteHeader(http.StatusNotFound)
	return
}


// CustomNotFoundHandler handles all requests not matched to a route.
func CustomNotFoundHandler (w http.ResponseWriter, r *http.Request) {
	// swagger:route GET /* NotFoundHandler NotFound
	// Works with all methods.
	//
	// Responses:
	// 404: NotFound
	w.WriteHeader(http.StatusNotFound)
}


// AddUser creates a new user with the id, name and password provided by client.
func AddUser (w http.ResponseWriter, r *http.Request) {
	// swagger:route POST /add AddUser AddUserParam
	//
	// responses:
	// 201: Created
	// 400: BadRequest
	// 401: Unauthorized
	// 409: Conflict
	user := &SingleUser{}
	auth:= r.Header.Get("X-Auth-Key")
	if auth != Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, exists := allUsers[user.ID]
	if exists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	allUsers[user.ID] = user

	w.WriteHeader(http.StatusCreated)
	return
}


// UpdateUser updates name, password or both for a user matched to id provided by client.
func UpdateUser (w http.ResponseWriter, r *http.Request) {
	// swagger:route PUT /update UpdateUser UserParam
	//
	// responses:
	// 204: NoContent
	// 400: BadRequest
	// 401: Unauthorized
	// 404: NotFound
	auth:= r.Header.Get("X-Auth-Key")
	if auth != Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user := &SingleUser{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := allUsers[user.ID]
	if ok {
		if user.Name != "" {
			allUsers[user.ID].Name = user.Name
		}

		if user.Password != "" {
			allUsers[user.ID].Password = user.Password
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}



// DeleteUser deletes a user matched to id provided by client.
func DeleteUser (w http.ResponseWriter, r *http.Request) {
	// swagger:route DELETE /delete DeleteUser deleteUser
	//
	// responses:
	// 204: NoContent
	// 400: BadRequest
	// 401: Unauthorized
	// 404: NotFound
	auth:= r.Header.Get("X-Auth-Key")
	if auth != Token {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := &SingleID{}
	err := json.NewDecoder(r.Body).Decode(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, exists := allUsers[id.ID]
	if exists {
		delete(allUsers, id.ID)

		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}


func main () {
	addr := "127.0.0.1:9000"
	mux := mux.NewRouter()

	server := http.Server{
		Handler: mux,
		Addr: addr,
		ReadTimeout: 3 * time.Second,
		WriteTimeout: 3 * time.Second,
		IdleTimeout: 10 * time.Second,
	}

	mux.NotFoundHandler = http.HandlerFunc(CustomNotFoundHandler)
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getMux := mux.Methods(http.MethodGet).Subrouter()
	getMux.HandleFunc("/get", GetSingleUser)
	getMux.Handle("/docs", sh)
	getMux.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	postMux := mux.Methods(http.MethodPost).Subrouter()
	postMux.HandleFunc("/add", AddUser)

	putMux := mux.Methods(http.MethodPut).Subrouter()
	putMux.HandleFunc("/update", UpdateUser)

	delMux := mux.Methods(http.MethodDelete).Subrouter()
	delMux.HandleFunc("/delete", DeleteUser)

	log.Println("Listening: ", addr)
	log.Fatal(server.ListenAndServe())
}
