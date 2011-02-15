package main

import(
	"os"
	"fmt"
	"physfs"
)

func main() {
	physfs.SetSaneConfig("test-go-physfs", "", "aoi", false, false)
	ver := physfs.GetVersion()
	fmt.Printf("Version:\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v\n\n", ver.Major, ver.Minor, ver.Patch)
	linkver := physfs.GetLinkedVersion()
	fmt.Printf("LinkedVersion:\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v\n\n", linkver.Major, linkver.Minor, linkver.Patch)

	//sat := physfs.SupportedArchiveTypes()
	//fmt.Printf("SupportedArchiveTypes: %v:\n", len(sat))
	//for i := range(sat) {
	//	fmt.Printf("\t%v: %v\n", i+1, sat[i].Extension)
	//}
	//fmt.Printf("\n")

	fmt.Printf("BaseDir: %v\n", physfs.GetBaseDir())
	fmt.Printf("UserDir: %v\n", physfs.GetUserDir())
	fmt.Printf("\n")

	buffer := make([]byte, 1024)
	file1, err := physfs.Open("dir1/file1", os.O_RDONLY)
	if err != nil {
		fmt.Printf("Error: %v\n", err.String())
		os.Exit(1)
	}
	defer file1.Close()
	n, _ := file1.Read(buffer)
	fmt.Printf("%v\n", string(buffer[0:n]))

	file2, err := physfs.Open("file2", os.O_WRONLY)
	if err != nil {
		fmt.Printf("Error: %v\n", err.String())
		os.Exit(1)
	}
	defer file2.Close()
	fmt.Fprintf(file2, "This is also a test.")

	sp, _ := physfs.GetSearchPath()
	fmt.Printf("%v\n\n", sp)

	list, _ := physfs.EnumerateFiles("/")
	fmt.Printf("%v\n\n", list)
}
