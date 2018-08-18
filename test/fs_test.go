package main

import (
	"os"
	"log"
	"testing"
	"../physfs"
)

func TestHTTPFileSystem(t *testing.T) {
	err := physfs.Init()
	if err != nil {
		log.Printf("0 Error: %v\n", err)
		t.FailNow()
	}
	physfs.Mount("a.zip", "", true)
	fs := physfs.FileSystem()

	t.Run("Stat() directory and IsDir()", func(t *testing.T) { 
		fi, err := fs.Open("/")
		if err != nil {
			log.Printf("Error cannot open: %v\n", err)
			t.FailNow()
		}
		info, err := fi.Stat()
		if err != nil {
			log.Printf("Error cannot stat: %v\n", err)
			t.FailNow()		
		}
		if !info.IsDir() {
			log.Println("not isdir")
			t.FailNow()
		}
	})
	t.Run("fs.Open FileSystem return IsNotExist", func(t *testing.T) { 
		_, err := fs.Open("/afasfasf")
		if !os.IsNotExist(err) {
			t.FailNow()
		}
	})
	
	err = physfs.Deinit()
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
}
