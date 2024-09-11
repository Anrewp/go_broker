package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Server struct {
	port string
	mux  *http.ServeMux
}

func (s *Server) Init() {
	s.port = fmt.Sprintf(":%v", *port)
	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", handle)
}

func (s *Server) Run() {
	log.Printf("Server listening on port %v\n\n", s.port)
	err := http.ListenAndServe(s.port, s.mux)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
