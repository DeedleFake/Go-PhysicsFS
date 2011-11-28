package main

import (
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"os"
)

func main() {
	err := physfs.Init()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	err = physfs.SetSaneConfig("test-go-physfs", "", "aoi", false, false)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	ver := physfs.VERSION()
	fmt.Printf("Version:\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v\n\n", ver.Major, ver.Minor, ver.Patch)
	linkver := physfs.GetLinkedVersion()
	fmt.Printf("LinkedVersion:\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v\n\n", linkver.Major, linkver.Minor, linkver.Patch)

	fmt.Printf("BaseDir: %v\n", physfs.GetBaseDir())
	fmt.Printf("UserDir: %v\n", physfs.GetUserDir())
	fmt.Printf("\n")

	buffer := make([]byte, 1024)
	file1, err := physfs.Open("dir1/file1")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	n, _ := file1.Read(buffer)
	fmt.Printf("%v\n", string(buffer[0:n]))
	file1.Close()

	file2, err := physfs.Create("file2")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(file2, "This is also a test.")
	file2.Close()

	sp, _ := physfs.GetSearchPath()
	fmt.Printf("%v\n\n", sp)

	physfs.GetSearchPathCallback(func(data interface{}, str string) {
		fmt.Printf("SearchPath From Callback: %v\n", str)
		fmt.Printf("Got Data: %v\n", data)
	}, "This is a test.")
	fmt.Printf("\n")

	physfs.GetCdRomDirsCallback(func(data interface{}, str string) {
		fmt.Printf("CdRomDirs From Callback: %v\n", str)
	}, nil)
	fmt.Printf("\n")

	list, _ := physfs.EnumerateFiles("/")
	fmt.Printf("%v\n\n", list)

	physfs.EnumerateFilesCallback("/", func(data interface{}, od string, fn string) {
		fmt.Printf("EnumerateFiles From Callback: %v/%v\n", od, fn)
	}, nil)
	fmt.Printf("\n")

	err = physfs.SetWriteDir(physfs.GetBaseDir())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Recursively deleting 'test-dir'...\n")
	err = physfs.DeleteRecurse("test-dir")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Succeeded.\n\n")
}
