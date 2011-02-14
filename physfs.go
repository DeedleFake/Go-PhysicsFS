package physfs

import(
	"os"
	"unsafe"
)

// #include <physfs.h>
import "C"

func init() {
	if C.PHYSFS_init(C.CString(os.Args[0])) == 0 {
		panic(C.PHYSFS_getLastError())
	}
}

type File C.PHYSFS_File

type ArchiveInfo struct {
	Extension string
	Description string
	Author string
	URL string
}

type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

func GetVersion() (ver *Version) {
	ver = new(Version)

	ver.Major = C.PHYSFS_VER_MAJOR
	ver.Minor = C.PHYSFS_VER_MINOR
	ver.Patch = C.PHYSFS_VER_PATCH

	return ver
}

func GetLinkedVersion() (ver *Version) {
	var v C.PHYSFS_Version
	C.PHYSFS_getLinkedVersion(&v)
	ver = (*Version)(unsafe.Pointer(&v))

	return ver
}

func SupportedArchiveTypes() (ai []*ArchiveInfo) {
	cai := C.PHYSFS_supportedArchiveTypes()

	i := uintptr(0)
	for {
		archive := (**C.PHYSFS_ArchiveInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cai)) + i)) //(*)
		if *archive == nil {
			break
		}

		ai = append(ai, (*ArchiveInfo)(unsafe.Pointer(*archive)))

		i += uintptr(unsafe.Sizeof(*archive))
	}

	return ai
}

func GetBaseDir() (string) {
	return C.GoString(C.PHYSFS_getBaseDir())
}

func GetUserDir() (string) {
	return C.GoString(C.PHYSFS_getUserDir())
}
