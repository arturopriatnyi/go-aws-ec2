package main

import (
	"log"
	"os"
)

var (
	l = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
)

func main() {
	l.Println("Hello from AWS EC2!")
}
