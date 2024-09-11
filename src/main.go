package main

import (
	"flag"
)

var port *int
var server Server
var queue = Queue{
	channels: make(map[string]chan string),
}

// using init is bad practice
func init() {
	port = flag.Int("port", 3000, "port on which server will be listened")
	flag.Parse()
}

func main() {
	server.Init()
	server.Run()
}
