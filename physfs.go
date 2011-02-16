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

// A type used to store information about supported archive types. Due to the
// fact that SupportedArchiveTypes() is currently broken, this is unused.
type ArchiveInfo struct {
	Extension string
	Description string
	Author string
	URL string
}

// Used to store information about the version PhysicsFS go-physfs was linked
// against.
type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// A type for functions to be called by the *Callback functions. Since they're
// currently broken, this is useless.
type StringCallback func(interface{}, string)

// Returns a boolean indicating if PhysicsFS has been initialized.
func IsInit() (bool) {
	if int(C.PHYSFS_isInit()) != 0 {
		return true
	}

	return false
}

// Deinitialize PhysicsFS. This closes any files that have been opened by
// PhysicsFS, clears the search and write paths, and cleans up other related
// resources.
func Deinit() (os.Error) {
	if int(C.PHYSFS_deinit()) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

// Returns a string containing an error message related to the last error
// that occured in a PhysicsFS function. Isn't necessary to call in most cases,
// as functions that generate said error return them as an os.Error in
// go-physfs.
func GetLastError() (string) {
	return C.GoString(C.PHYSFS_getLastError())
}

// Returns a Version containing the version of PhysicsFS that the bindings were
// compiled against.
func VERSION() (ver *Version) {
	ver = new(Version)

	ver.Major = C.PHYSFS_VER_MAJOR
	ver.Minor = C.PHYSFS_VER_MINOR
	ver.Patch = C.PHYSFS_VER_PATCH

	return ver
}

// Returns a Version containing the version of PhysicsFS that the bindings are
// linked against.
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

// Returns the the directory in which the application is. May or may not
// correspond to the processes current working directory.
func GetBaseDir() (string) {
	return C.GoString(C.PHYSFS_getBaseDir())
}

// Returns the home directory of the user that ran the application.
func GetUserDir() (string) {
	return C.GoString(C.PHYSFS_getUserDir())
}

// Returns the current write directory. Files written using PhysicsFS can only
// be inside the write directory. Default is nowhere, which will return a blank
// string.
func GetWriteDir() (string) {
	return C.GoString(C.PHYSFS_getWriteDir())
}

// Set the current write directory. Returns an error, if any.
func SetWriteDir(dir string) (os.Error) {
	if int(C.PHYSFS_setWriteDir(C.CString(dir))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

// Gets the directory seperator for the operating system. In Windows returns
// "\\", in Linux "/", and in MacOS versions before OS X returns ":".
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
