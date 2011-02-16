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

type StringCallback func(interface{}, string)

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

func GetDirSeparator() (string) {
	return C.GoString(C.PHYSFS_getDirSeparator())
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

func GetCdRomDirs() (sp []string, err os.Error) {
	csp := C.PHYSFS_getCdRomDirs()

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

	C.PHYSFS_freeList(unsafe.Pointer(csp))
	return sp, nil
}

//func GetCdRomDirsCallback(c StringCallback, d interface{}) {
//	C.PHYSFS_getCdRomDirsCallback((*[0]uint8)(unsafe.Pointer(&c)), unsafe.Pointer(&d))
//}

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

	C.PHYSFS_freeList(unsafe.Pointer(csp))
	return sp, nil
}

//func GetSearchPathCallback(c StringCallback, d interface{}) {
//	C.PHYSFS_getSearchPathCallback((*[0]uint8)(unsafe.Pointer(&c)), unsafe.Pointer(&d))
//}

func PermitSymbolicLinks(set bool) {
	s := C.int(0)
	if set {
		s = 1
	}

	C.PHYSFS_permitSymbolicLinks(s)
}

func SymbolicLinksPermitted() (bool) {
	if int(C.PHYSFS_symbolicLinksPermitted()) != 0 {
		return true
	}

	return false
}

func IsSymbolicLink(n string) (bool) {
	if int(C.PHYSFS_isSymbolicLink(C.CString(n))) != 0 {
		return true
	}

	return false
}

func GetRealDir(n string) (string, os.Error) {
	dir := C.PHYSFS_getRealDir(C.CString(n))

	if dir != nil {
		return C.GoString(dir), nil
	}

	return C.GoString(dir), os.NewError(GetLastError())
}

func EnumerateFiles(dir string) (list []string, err os.Error) {
	clist := C.PHYSFS_enumerateFiles(C.CString(dir))

	if clist == nil {
		return nil, os.NewError(GetLastError())
	}

	i := uintptr(0)
	for {
		item := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(clist)) + i))
		if item == nil {
			break
		}

		list = append(list, C.GoString(item))

		i += uintptr(unsafe.Sizeof(clist))
	}

	C.PHYSFS_freeList(unsafe.Pointer(clist))
	return list, nil
}

func Exists(n string) (bool) {
	if int(C.PHYSFS_exists(C.CString(n))) != 0 {
		return true
	}

	return false
}

func Delete(n string) (os.Error) {
	if int(C.PHYSFS_delete(C.CString(n))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func IsDirectory(dir string) (bool) {
	if int(C.PHYSFS_isDirectory(C.CString(dir))) != 0 {
		return true
	}

	return false
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

func GetMountPoint(dir string) (string, os.Error) {
	mp := C.PHYSFS_getMountPoint(C.CString(dir))

	if mp != nil {
		return C.GoString(mp), nil
	}

	return C.GoString(mp), os.NewError(GetLastError())
}

func AddToSearchPath(dir string, app bool) (os.Error) {
	a := 0
	if app {
		a = 1
	}

	if int(C.PHYSFS_addToSearchPath(C.CString(dir), C.int(a))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func RemoveFromSearchPath(dir string) (os.Error) {
	if int(C.PHYSFS_removeFromSearchPath(C.CString(dir))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func GetLastModTime(n string) (int64, os.Error) {
	num := int64(C.PHYSFS_getLastModTime(C.CString(n)))

	if num != -1 {
		return num, nil
	}

	return num, os.NewError(GetLastError())
}
