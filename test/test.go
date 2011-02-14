package main

import(
	"os"
	"fmt"
	"physfs"
)

func main() {
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
	physfs.Mount("zip1.zip", "", true)
	file1 := physfs.Open("dir1/file1", os.O_RDONLY)
	if file1 == nil {
		fmt.Printf("Error: %v\n", physfs.GetLastError())
		os.Exit(1)
	}
	defer file1.Close()
	n := file1.Read(buffer)
	fmt.Printf("%v", string(buffer[0:n]))
}
