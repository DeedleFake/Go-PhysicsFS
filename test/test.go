package main

import(
	"fmt"
	"physfs"
)

func main() {
	fmt.Printf("BaseDir: %v\n", physfs.GetBaseDir())
	fmt.Printf("UserDir: %v\n", physfs.GetUserDir())
}
