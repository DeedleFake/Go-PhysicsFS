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

func Deinit() (os.Error) {
	if int(C.PHYSFS_deinit()) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
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

//func SupportedArchiveTypes() (ai []ArchiveInfo) {
//	cai := C.PHYSFS_supportedArchiveTypes()
//
//	i := uintptr(0)
//	for {
//		archive := *(**C.PHYSFS_ArchiveInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cai)) + i))
//		if archive == nil {
//			break
//		}
//
//		ai = append(ai, *(*ArchiveInfo)(unsafe.Pointer(archive)))
//
//		i += uintptr(unsafe.Sizeof(cai))
//	}
//
//	return ai
//}

func GetBaseDir() (string) {
	return C.GoString(C.PHYSFS_getBaseDir())
}

func GetUserDir() (string) {
	return C.GoString(C.PHYSFS_getUserDir())
}

func GetWriteDir() (string) {
	return C.GoString(C.PHYSFS_getWriteDir())
}

func SetWriteDir(dir string) (os.Error) {
	if int(C.PHYSFS_setWriteDir(C.CString(dir))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func SetSaneConfig(org, app, ext string, cd, arc bool) (os.Error) {
	cdArg := 0
	if cd {
		cdArg = 1
	}

	arcArg := 0
	if arc {
		arcArg = 1
	}

	if int(C.PHYSFS_setSaneConfig(C.CString(org), C.CString(app), C.CString(ext), C.int(cdArg), C.int(arcArg))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func GetSearchPath() (sp []string, err os.Error) {
	csp := C.PHYSFS_getSearchPath()

	if csp == nil {
		return nil, os.NewError(GetLastError())
	}

	i := uintptr(0)
	for {
		p := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(csp)) + i))
		if p == nil {
			break
		}

		sp = append(sp, C.GoString(p))

		i += uintptr(unsafe.Sizeof(csp))
	}

	return sp, nil
}

func PermitSymbolicLinks(set bool) {
	s := C.int(0)
	if set {
		s = 1
	}

	C.PHYSFS_permitSymbolicLinks(s)
}

func Mkdir(dir string) (os.Error) {
	if int(C.PHYSFS_mkdir(C.CString(dir))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func Mount(dir, mp string, app bool) (os.Error) {
	a := 0
	if app {
		a = 1
	}

	if int(C.PHYSFS_mount(C.CString(dir), C.CString(mp), C.int(a))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func Open(name string, flag int) (f *File, err os.Error) {
	switch flag {
		case os.O_RDONLY:
			f = (*File)(C.PHYSFS_openRead(C.CString(name)))
		case os.O_WRONLY:
			f = (*File)(C.PHYSFS_openWrite(C.CString(name)))
		case os.O_WRONLY | os.O_APPEND:
			f = (*File)(C.PHYSFS_openAppend(C.CString(name)))
		default:
			return nil, os.NewError("Unknown flag(s).")
	}

	if f == nil {
		return f, os.NewError(GetLastError())
	}

	return f, nil
}

func (f *File)Close() (os.Error) {
	if int(C.PHYSFS_close((*C.PHYSFS_File)(f))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func (f *File)Read(buf []byte) (n int, err os.Error) {
	n = int(C.PHYSFS_read((*C.PHYSFS_File)(f), unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))

	if n == -1 {
		return n, os.NewError(GetLastError())
	}

	return n, nil
}

func (f *File)Write(buf []byte) (n int, err os.Error) {
	n = int(C.PHYSFS_write((*C.PHYSFS_File)(f), unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))

	if n == -1 {
		return n, os.NewError(GetLastError())
	}

	return n, nil
}

func (f *File)EOF() (bool) {
	if int(C.PHYSFS_eof((*C.PHYSFS_File)(f))) != 0 {
		return true
	}

	return false
}

func (f *File)Tell() (int64, os.Error) {
	r := int64(C.PHYSFS_tell((*C.PHYSFS_File)(f)))
	if r == -1 {
		return r, os.NewError(GetLastError())
	}

	return r, nil
}

func (f *File)Seek(offset int64, none int) (int64, os.Error) {
	r := int64(C.PHYSFS_seek((*C.PHYSFS_File)(f), C.PHYSFS_uint64(offset)))

	if r == 0 {
		return r, os.NewError(GetLastError())
	}

	return r, nil
}
