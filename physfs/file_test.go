package physfs

import(
	"os"
	"fmt"
	"testing"
)

func TestFile(t *testing.T) {
	if !IsInit() {
		err := Init()
		if err != nil {
			t.Fatalf("Error: %v\n", err.String())
		}
	}

	err := Mount("../test/zip1.aoi", "", true)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}

	err = Mount("../test", "dir2", true)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}

	file1, err := Open("dir1/file1", os.O_RDONLY)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
	buffer := make([]byte, 256)
	n, err := file1.Read(buffer)
	if (err != nil) && (err != os.EOF) {
		t.Fatalf("Error: %v\n", err.String())
	}
	fmt.Printf("%v\n", string(buffer[0:n]))
	file1.Close()

	SetWriteDir("../test")
	file2, err := Open("file2", os.O_WRONLY)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
	fmt.Fprintf(file2, "This is a test.\nThis is also a test.\nThis is yet another test.")
	file2.Close()

	file2, err = Open("dir2/file2", os.O_RDONLY)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
	_, err = file2.Seek(-60, 2)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
	np, err := file2.Seek(14, 1)
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
	n, err = file2.Read(buffer)
	if (err != nil) && (err != os.EOF) {
		t.Fatalf("Error: %v\n", err.String())
	}
	fmt.Printf("%v: %v\n", np, string(buffer[0:n]))
	file2.Close()

	err = Delete("file2")
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}

	err = Deinit()
	if err != nil {
		t.Fatalf("Error: %v\n", err.String())
	}
}
