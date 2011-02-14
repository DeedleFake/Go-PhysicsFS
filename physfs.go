package physfs

import(
	"os"
)

// #include <physfs.h>
import "C"

func init() {
	if C.PHYSFS_init(C.CString(os.Args[0])) == 0 {
		panic(C.PHYSFS_getLastError())
	}
}

func GetBaseDir() (string) {
	return C.GoString(C.PHYSFS_getBaseDir())
}
