package physfs

import(
	"os"
	"unsafe"
)

// #include <physfs.h>
import "C"

func init() {
	if !IsInit() {
		if C.PHYSFS_init(C.CString(os.Args[0])) == 0 {
			panic(GetLastError())
		}
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

func IsInit() (bool) {
	if int(C.PHYSFS_isInit()) != 0 {
		return true
	}

	return false
}

func Deinit() (int) {
	return int(C.PHYSFS_deinit())
}

func GetLastError() (string) {
	return C.GoString(C.PHYSFS_getLastError())
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

func Mount(dir, mp string, app bool) (bool) {
	a := 0
	if app {
		a = 1
	}

	if int(C.PHYSFS_mount(C.CString(dir), C.CString(mp), C.int(a))) != 0 {
		return true
	}

	return false
}

func Open(name string, flag int) (*File) {
	switch flag {
		case os.O_RDONLY:
			return (*File)(C.PHYSFS_openRead(C.CString(name)))
		default:
			panic("Unknown flag.")
	}

	panic("What the heck?!?")
	return nil
}

func (f *File)Close() (bool) {
	if int(C.PHYSFS_close((*C.PHYSFS_File)(f))) != 0 {
		return true
	}

	return false
}

func (f *File)Read(buf []byte) (int) {
	return int(C.PHYSFS_read((*C.PHYSFS_File)(f), unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))
}
