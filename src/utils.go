package main

import (
	"log"
	"os"
)

func checkEnv() {
	_, ok := os.LookupEnv(PORT)
	if !ok {
		log.Println("Set PORT env before run")
		os.Exit(1)
	}
}
