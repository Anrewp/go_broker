package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const MSG = "v"
const TIMEOUT = "timeout"

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
		badRequest(w, qName)
		return
	}
	queue.Push(qName, msg)
	succesPush(w, qName, msg)
}

func handleGET(w http.ResponseWriter, r *http.Request) {
	timeout := r.URL.Query().Get(TIMEOUT)
	if timeout != "" {
		handleWithTimeout(w, r, timeout)
	} else {
		handleCommon(w, r)
	}
}

func handleCommon(w http.ResponseWriter, r *http.Request) {
	qName := rQueue(r)
	msg := queue.Pop(qName)
	if msg == "" {
		notFound(w, qName)
		return
	}

	succesPop(w, qName, msg)
}

func handleWithTimeout(w http.ResponseWriter, r *http.Request, timeout string) {
	qName := rQueue(r)
	ttl, err := strconv.Atoi(timeout)
	if err != nil {
		log.Println(err)
		return
	}
	msg := queue.PopWait(qName, time.Duration(ttl))
	if msg == "" {
		notFound(w, qName)
		return
	}

	succesPop(w, qName, msg)
}

func succesPush(w http.ResponseWriter, qName, msg string) {
	w.WriteHeader(http.StatusOK)
	log.Printf("200 OK [%v]: <- %v ", qName, msg)
}

func succesPop(w http.ResponseWriter, qName, msg string) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
	log.Printf("200 OK [%v]: -> %v ", qName, msg)
}

func notFound(w http.ResponseWriter, qName string) {
	w.WriteHeader(http.StatusNotFound)
	log.Printf("404 NotFound [%v]: -> nil", qName)
}

func badRequest(w http.ResponseWriter, qName string) {
	w.WriteHeader(http.StatusBadRequest)
	log.Printf("400 BadRequest [%v]: <- nil", qName)
}

func rQueue(r *http.Request) string {
	return strings.Split(r.URL.Path, "/")[1]
}

func rMsg(r *http.Request) string {
	return r.URL.Query().Get(MSG)
}
