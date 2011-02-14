package main

import(
	"fmt"
	"physfs"
)

func main() {
	ver := physfs.GetVersion()
	fmt.Printf("Version:\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v\n\n", ver.Major, ver.Minor, ver.Patch)
	linkver := physfs.GetLinkedVersion()
	fmt.Printf("LinkedVersion:\n\tMajor: %v\n\tMinor: %v\n\tPatch: %v\n\n", linkver.Major, linkver.Minor, linkver.Patch)

	sat := physfs.SupportedArchiveTypes()
	fmt.Printf("SupportedArchiveTypes: %v:\n", len(sat))
	for i := range(sat) {
		fmt.Printf("\t%v: %v\n", i+1, sat[i].Extension)
	}
	fmt.Printf("\n")

	fmt.Printf("BaseDir: %v\n", physfs.GetBaseDir())
	fmt.Printf("UserDir: %v\n", physfs.GetUserDir())
}
