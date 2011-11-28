// The physfs package provides Go bindings for the PhysicsFS
// archive-abstraction library.
package physfs

import (
	"errors"
	"os"
	"path"
	"time"
	"unsafe"
)

// #cgo LDFLAGS: -lphysfs
// #include <stdlib.h>
// #include <physfs.h>
//
// #include "wrapcb.h"
import "C"

const (
	T_LOCAL = iota
	T_UTC
)

// A type used to store information about supported archive types.
type ArchiveInfo struct {
	Extension   string
	Description string
	Author      string
	URL         string
}

// Used to store information about the version of PhysicsFS go-physfs was linked
// against.
type Version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// The func type required by GetCdRomDirsCallback and GetSearchPathCallback.
type StringCallback func(interface{}, string)

type contextStringCallback struct {
	cb   StringCallback
	data interface{}
}

//export wrapStringCallback
func wrapStringCallback(data unsafe.Pointer, str *C.char) {
	csc := (*contextStringCallback)(data)
	csc.cb(csc.data, C.GoString(str))
}

// The func type required by EnumerateFilesCallback.
type EnumFilesCallback func(interface{}, string, string)

type contextEnumFilesCallback struct {
	cb   EnumFilesCallback
	data interface{}
}

//export wrapEnumFilesCallback
func wrapEnumFilesCallback(data unsafe.Pointer, origdir *C.char, fname *C.char) {
	cefc := (*contextEnumFilesCallback)(data)
	cefc.cb(cefc.data, C.GoString(origdir), C.GoString(fname))
}

// Returns a boolean indicating if PhysicsFS has been initialized.
func IsInit() bool {
	if int(C.PHYSFS_isInit()) != 0 {
		return true
	}

	return false
}

// Initialize PhysicsFS. Must be called before most functions will work. Returns
// an error, if any.
func Init() error {
	arg0 := C.CString(os.Args[0])
	defer C.free(unsafe.Pointer(arg0))
	if int(C.PHYSFS_init(arg0)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Deinitialize PhysicsFS. This closes any files that have been opened by
// PhysicsFS, clears the search and write paths, forgets other settings, such as
// whether or not symbolic links are permitted, and cleans up other related
// resources.
func Deinit() error {
	if int(C.PHYSFS_deinit()) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Returns a string containing an error message related to the last error
// that occured in a PhysicsFS function. Isn't necessary to call in most cases,
// as functions that generate said error return them as an error in
// go-physfs.
func GetLastError() string {
	cerr := C.PHYSFS_getLastError()
	return C.GoString(cerr)
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

// Returns an []ArchiveInfo containing information about all the archives
// supported by PhysicsFS.
func SupportedArchiveTypes() (ai []ArchiveInfo) {
	cai := C.PHYSFS_supportedArchiveTypes()

	i := uintptr(0)
	for {
		archive := *(**C.PHYSFS_ArchiveInfo)(unsafe.Pointer(uintptr(unsafe.Pointer(cai)) + i))
		if archive == nil {
			break
		}

		var a ArchiveInfo
		a.Extension = C.GoString(archive.extension)
		a.Description = C.GoString(archive.description)
		a.Author = C.GoString(archive.author)
		a.URL = C.GoString(archive.url)

		ai = append(ai, a)

		i += uintptr(unsafe.Sizeof(cai))
	}

	return ai
}

// Returns the the directory in which the application is. May or may not
// correspond to the processes current working directory.
func GetBaseDir() string {
	cdir := C.PHYSFS_getBaseDir()
	return C.GoString(cdir)
}

// Returns the home directory of the user that ran the application.
func GetUserDir() string {
	cdir := C.PHYSFS_getUserDir()
	return C.GoString(cdir)
}

// Returns the current write directory. Files written using PhysicsFS can only
// be inside the write directory. Default is nowhere, which will return a blank
// string.
func GetWriteDir() string {
	cdir := C.PHYSFS_getWriteDir()
	return C.GoString(cdir)
}

// Set the current write directory. Returns an error, if any.
func SetWriteDir(dir string) error {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	if int(C.PHYSFS_setWriteDir(cdir)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Gets the directory seperator for the operating system. In Windows returns
// "\\", in Linux "/", and in MacOS versions before OS X returns ":".
func GetDirSeparator() string {
	cdir := C.PHYSFS_getDirSeparator()
	return C.GoString(cdir)
}

// Sets up sane, default search/write paths. The write path is set to
// 'GetUserDir()/.org/app', which is created if it doesn't exist. The search
// path is set to the write path, GetBaseDir(), any deteced CD-ROM directories,
// if specified by cd, and any archives found in any of the previously listed
// locations that have extensions that match ext. To not automatically load
// archives, simply give a blank string. Do not specifiy a '.' before the
// extension. If pre is true the archives are prepended to the search path; if
// false they are appended.
func SetSaneConfig(org, app, ext string, cd, pre bool) error {
	cdArg := 0
	if cd {
		cdArg = 1
	}

	preArg := 0
	if pre {
		preArg = 1
	}

	corg := C.CString(org)
	defer C.free(unsafe.Pointer(corg))
	capp := C.CString(app)
	defer C.free(unsafe.Pointer(capp))
	cext := C.CString(ext)
	defer C.free(unsafe.Pointer(cext))

	r := C.PHYSFS_setSaneConfig(corg, capp, cext, C.int(cdArg), C.int(preArg))
	if int(r) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Returns a []string containing all detected CD-ROM directories and an error,
// if any. Note that detection of CD-ROM drives is dependent on various factors,
// such as whether or not there is a disc in the drive. Also note that while
// this function and related ones refer to CD-ROMs, they will detect any type of
// supported disc, including DVDs and Blu-Ray discs.
func GetCdRomDirs() (sp []string, err error) {
	csp := C.PHYSFS_getCdRomDirs()

	if csp == nil {
		return nil, errors.New(GetLastError())
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

// Call c for each detected CD-ROM directrory, passing it d and the dir.
func GetCdRomDirsCallback(c StringCallback, d interface{}) {
	csc := &contextStringCallback{
		cb:   c,
		data: d,
	}

	C.getCdRomDirsCallback(unsafe.Pointer(csc))
}

// Call c for each entry in the SearchPath, passing it d and the dir.
func GetSearchPathCallback(c StringCallback, d interface{}) {
	csc := &contextStringCallback{
		cb:   c,
		data: d,
	}

	C.getSearchPathCallback(unsafe.Pointer(csc))
}

// Call c for each file in dir, passing it d, dir, and the file.
func EnumerateFilesCallback(dir string, c EnumFilesCallback, d interface{}) {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))

	cefc := &contextEnumFilesCallback{
		cb:   c,
		data: d,
	}

	C.enumerateFilesCallback(cdir, unsafe.Pointer(cefc))
}

// Returns a []string with the current search path, in order, and an error, if
// any.
func GetSearchPath() (sp []string, err error) {
	csp := C.PHYSFS_getSearchPath()

	if csp == nil {
		return nil, errors.New(GetLastError())
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

// Enable or disable the following of symbolic links. Default is disabled.
func PermitSymbolicLinks(set bool) {
	s := C.int(0)
	if set {
		s = 1
	}

	C.PHYSFS_permitSymbolicLinks(s)
}

// Return whether or not following of symbolic links is currently enabled.
func SymbolicLinksPermitted() bool {
	if int(C.PHYSFS_symbolicLinksPermitted()) != 0 {
		return true
	}

	return false
}

// Returns true if the named path exists and is a symbolic link, otherwise
// returns false.
func IsSymbolicLink(n string) bool {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	if int(C.PHYSFS_isSymbolicLink(cn)) != 0 {
		return true
	}

	return false
}

// Returns the real path to the specified file/directory. For example, if you
// call:
//		physfs.GetRealDir("maps/level1.map")
// and 'maps/level1.map' actually exists at 'C:\mygame\maps\level1.map', and
// 'C:\mygame' is in your search path, 'C:\mygame' is returned. Also returns an
// error, if any.
func GetRealDir(n string) (string, error) {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	dir := C.PHYSFS_getRealDir(cn)

	if dir != nil {
		return C.GoString(dir), nil
	}

	return C.GoString(dir), errors.New(GetLastError())
}

// Returns a []string containing the files and directories in the specified
// directory in your search path, and an error, if any.
func EnumerateFiles(dir string) (list []string, err error) {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	clist := C.PHYSFS_enumerateFiles(cdir)

	if clist == nil {
		return nil, errors.New(GetLastError())
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

// Returns a boolean indicating whether or not the specified file/directory
// exists.
func Exists(n string) bool {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	if int(C.PHYSFS_exists(cn)) != 0 {
		return true
	}

	return false
}

// Deletes the specified file or directory. Only deletes empty directories.
// Returns an error, if any.
func Delete(n string) error {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	if int(C.PHYSFS_delete(cn)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// A convienece function that will recurse through a directory, deleting all
// sub-directories and files and then the original directory as well. Returns an
// error, if any.
func DeleteRecurse(dir string) (err error) {
	if IsDirectory(dir) {
		files, err := EnumerateFiles(dir)
		if err != nil {
			return err
		}
		for _, i := range files {
			diri := path.Join(dir, i)
			if IsDirectory(diri) {
				err = DeleteRecurse(diri)
				if err != nil {
					return err
				}
			} else if Exists(diri) {
				err = Delete(diri)
				if err != nil {
					return err
				}
			}
		}

		err = Delete(dir)
		if err != nil {
			return err
		}
	}

	return err
}

// Returns true if dir exists and is a directory. Otherwise, returns false.
func IsDirectory(dir string) bool {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	if int(C.PHYSFS_isDirectory(cdir)) != 0 {
		return true
	}

	return false
}

// Creates the specified directory inside the write path. Will create any parent
// directories that don't exist. Returns an error, if any.
func Mkdir(dir string) error {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	if int(C.PHYSFS_mkdir(cdir)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Adds an archive or directory dir to the search path, mounting it at the
// specified point mp in the search path. Use "" or "/" to simply add it to the
// search path. If app is true dir is appended to the search path; otherwise it
// is prepended. While multiple archives/directories may be mounted on the same
// mount-point, you may not mount the same archive/directory in multiple
// locations. Attempting to do so will simply do nothing without returning an
// error. Returns an error, if any.
func Mount(dir, mp string, app bool) error {
	a := 0
	if app {
		a = 1
	}

	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	cmp := C.CString(mp)
	defer C.free(unsafe.Pointer(cmp))

	if int(C.PHYSFS_mount(cdir, cmp, C.int(a))) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Gets the mount-point of the specified archive/directory. Returns the
// mount-point (Big surprise...) and an error, if any.
func GetMountPoint(dir string) (string, error) {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	mp := C.PHYSFS_getMountPoint(cdir)

	if mp != nil {
		return C.GoString(mp), nil
	}

	return C.GoString(mp), errors.New(GetLastError())
}

// A legacy function that is now equivalent to
//		physfs.Mount(dir, "", app)
func AddToSearchPath(dir string, app bool) error {
	a := 0
	if app {
		a = 1
	}

	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))

	if int(C.PHYSFS_addToSearchPath(cdir, C.int(a))) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Remove the specified archive/directory from search path. This will fail if
// there any files inside the archive/directory that are still open. Returns an
// error, if any.
func RemoveFromSearchPath(dir string) error {
	cdir := C.CString(dir)
	defer C.free(unsafe.Pointer(cdir))
	if int(C.PHYSFS_removeFromSearchPath(cdir)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Returns the last time the specified file was modified in either or the local
// time-zone or UTC, and an error, if any.
func GetLastModTime(n string, zone int) (t *time.Time, err error) {
	cn := C.CString(n)
	defer C.free(unsafe.Pointer(cn))
	num := int64(C.PHYSFS_getLastModTime(cn))

	if num != -1 {
		err = nil
		switch zone {
		case T_LOCAL:
			t = time.SecondsToLocalTime(num)
		case T_UTC:
			t = time.SecondsToUTC(num)
		default:
			err = errors.New("Unknown zone.")
		}

		return t, err
	}

	return t, errors.New(GetLastError())
}
