package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var server Server

type Server struct {
	port string
	mux  *http.ServeMux
}

func (s *Server) Init() {
	checkEnv()
	s.port = fmt.Sprintf(":%v", os.Getenv(PORT))
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", handle)
}

func (s *Server) Run() {
	err := http.ListenAndServe(s.port, s.mux)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
