package main

import (
	"log"
	"http"
	"physfs"
)

func main() {
	err := physfs.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer physfs.Deinit()

	err = physfs.Mount("test.zip", "/", true)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(physfs.FileSystem()))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
