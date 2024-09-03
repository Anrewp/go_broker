package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGET(w, r)
	case http.MethodPut:
		handlePUT(w, r)
	default:
		http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
	}
}

func handlePUT(w http.ResponseWriter, r *http.Request) {
	qName := rQueue(r)
	msg := rMsg(r)
	if msg == "" {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("400 BadRequest")
		return
	}

	queue.Push(qName, msg)
	w.WriteHeader(http.StatusOK)
	log.Println("200 OK pushed " + msg + " to queue " + qName)
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get(TIMEOUT)
	if timeout != "" {
		handleWithTimeOut(w, r, timeout)
	} else {
		handleCommon(w, r)
	}
}

func handleCommon(w http.ResponseWriter, r *http.Request) {
	qName := rQueue(r)
	msg := queue.Pop(qName)
	if msg == "" {
		w.WriteHeader(http.StatusNotFound)
		log.Println("404 Not Found")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
	log.Println("200 OK returned: " + msg + " from queue " + qName)
}

func handleWithTimeOut(w http.ResponseWriter, r *http.Request, timout string) {
	qName := rQueue(r)
	ttl, err := strconv.Atoi(timout)
	if err != nil {
		log.Println(err)
		return
	}
	ch := queue.PopWait(qName)

	select {
	case msg := <-ch:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
		log.Println("200 OK returned: " + msg + " from queue " + qName)
	case <-time.After(time.Duration(ttl) * time.Second):
		queue.RemoveSub(qName, ch)
		w.WriteHeader(http.StatusNotFound)
		log.Println("404 Not found with timeout")
	}
}

func rQueue(r *http.Request) string {
	return strings.Split(r.URL.Path, "/")[1]
}

func rMsg(r *http.Request) string {
	return r.URL.Query().Get(MSG)
}
