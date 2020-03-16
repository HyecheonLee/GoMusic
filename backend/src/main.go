package main

import (
	"backend/src/rest"
	"log"
)

func main() {
	log.Println("Main log....")
	rest.RunAPI("localhost:8080")
}
