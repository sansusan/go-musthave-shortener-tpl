package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/oklog/ulid/v2"
	"io"
	"log"
	"net/http"
	"sync"
)

var store = make(map[string]string)
var mutex = sync.RWMutex{}

func main() {
	StartService()
}

func StartService() {
	router := httprouter.New()
	router.POST("/", ShortUrl)
	router.GET("/:id", GetUrlById)

	fmt.Printf("Starting server\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func ShortUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	url, _ := io.ReadAll(r.Body)
	mutex.Lock()
	defer mutex.Unlock()
	var uid = ulid.Make().String()
	store[uid] = string(url)
	w.WriteHeader(201)
	fmt.Fprintf(w, uid)
}

func GetUrlById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	mutex.RLock()
	defer mutex.RUnlock()
	var url, exist = store[(ps.ByName("id"))]
	if exist {
		w.Header().Set("Location", url)
		w.WriteHeader(307)
		fmt.Fprintf(w, url)
	} else {
		w.WriteHeader(404)
	}
}
