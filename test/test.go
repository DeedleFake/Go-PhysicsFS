package main

import(
	"fmt"
	"physfs"
)

func main() {
	fmt.Printf("%v\n", physfs.GetBaseDir())
}
