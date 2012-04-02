package physfs

import (
	"fmt"
	"io"
	"testing"
)

func TestFile(t *testing.T) {
	if !IsInit() {
		err := Init()
		if err != nil {
			t.Fatalf("Error: %v\n", err.Error())
		}
	}

	err := Mount("../test/zip1.aoi", "", true)
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}

	err = Mount("../test", "dir2", true)
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}

	file1, err := Open("dir1/file1")
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}
	buffer := make([]byte, 256)
	n, err := file1.Read(buffer)
	if (err != nil) && (err != io.EOF) {
		t.Fatalf("Error: %v\n", err.Error())
	}
	fmt.Printf("%v\n", string(buffer[0:n]))
	file1.Close()

	SetWriteDir("../test")
	file2, err := Create("file2")
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}
	fmt.Fprintf(file2, "This is a test.\nThis is also a test.\nThis is yet another test.")
	file2.Close()

	file2, err = Open("dir2/file2")
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}
	_, err = file2.Seek(-60, 2)
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}
	np, err := file2.Seek(14, 1)
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}
	n, err = file2.Read(buffer)
	if (err != nil) && (err != io.EOF) {
		t.Fatalf("Error: %v\n", err.Error())
	}
	fmt.Printf("%v: %v\n", np, string(buffer[0:n]))
	file2.Close()

	err = Delete("file2")
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}

	dir, err := Open("/")
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
	defer dir.Close()
	fi, err := dir.Readdir(-1)
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}
	fmt.Println("Readdir() Test:")
	for i := range fi {
		fmt.Printf("\t%v\n", fi[i].Name())
	}

	err = Deinit()
	if err != nil {
		t.Fatalf("Error: %v\n", err.Error())
	}
}
