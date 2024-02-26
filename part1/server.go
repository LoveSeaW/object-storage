package main

import (
	"log"
	"net/http"
	"object-storage/part1/objects"
)

func main() {
	http.HandleFunc("/object/", objects.Handler)
	log.Fatal(http.ListenAndServe(":8800", nil))
}
