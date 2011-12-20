package physfs

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"unsafe"
)

// #include <stdlib.h>
// #include <physfs.h>
import "C"

// A type for PhysicsFS file operations. Designed to be as compatible as
// possible with os.File.
type File struct {
	cfile *C.PHYSFS_File

	name string
	read int
}

// Open the named file, relative to the current write dir, for writing. The
// specified file is created if it doesn't exist. If it does exist it is
// truncated to zero bytes. Returns the file and an error, if any.
func Create(name string) (file *File, err error) {
	return openFile(name, os.O_WRONLY)
}

// Open the named file, relative the current write dir, for writing. The
// specified file is created if it doesn't exist. If it does exist, the writing
// offset is set to the end of the file so that writes will append to the file.
// Returns the file and an error, if any.
func Append(name string) (file *File, err error) {
	return openFile(name, os.O_APPEND)
}

// Open the named file from the search path for reading. Returns the file and an
// error, if any.
func Open(name string) (file *File, err error) {
	return openFile(name, os.O_RDONLY)
}

// Open the named file with the mode specified by flag. Read-only files are
// opened from the search-path, and write-only and append-only files are opened
// relative to the current write directory. Returns the file and an error, if
// any.
func openFile(name string, flag int) (f *File, err error) {
	if IsDirectory(name) {
		return &File{
			nil,
			name,
			0,
		}, nil
	}

	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	switch flag {
	case os.O_RDONLY:
		f = &File{C.PHYSFS_openRead(cname), name, -1}
	case os.O_WRONLY:
		f = &File{C.PHYSFS_openWrite(cname), name, -1}
	case os.O_APPEND:
		f = &File{C.PHYSFS_openAppend(cname), name, -1}
	default:
		return nil, errors.New("Unknown flag.")
	}

	if f.cfile == nil {
		return nil, errors.New(GetLastError())
	}

	return
}

func (f *File) isdir() bool {
	return IsDirectory(f.name)
}

// Close the file, release related resources. Returns an error, if any.
func (f *File) Close() error {
	if f.isdir() {
		return nil
	}

	if int(C.PHYSFS_close(f.cfile)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Read up to len(buf) bytes from the file into buf. Returns the number of bytes
// read and an error, if any.
func (f *File) Read(buf []byte) (n int, err error) {
	if f.isdir() {
		return 0, os.EISDIR
	}

	n = int(C.PHYSFS_read(f.cfile, unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))

	if n == -1 {
		err = errors.New(GetLastError())
	}

	if f.EOF() {
		err = io.EOF
	}

	return n, err
}

// Write the bytes in buf to the file. Returns the number of bytes written and
// an error, if any.
func (f *File) Write(buf []byte) (n int, err error) {
	if f.isdir() {
		return 0, os.EISDIR
	}

	n = int(C.PHYSFS_write(f.cfile, unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))

	if n == -1 {
		return n, errors.New(GetLastError())
	}

	return n, nil
}

// Returns a boolean indicating whether or not the end of the file has been
// reached.
func (f *File) EOF() bool {
	if f.isdir() {
		return true
	}

	if int(C.PHYSFS_eof(f.cfile)) != 0 {
		return true
	}

	return false
}

// Returns a number indication the current position in the file, and an error,
// if any.
func (f *File) Tell() (int64, error) {
	if f.isdir() {
		return 0, os.EISDIR
	}

	r := int64(C.PHYSFS_tell(f.cfile))
	if r == -1 {
		return r, errors.New(GetLastError())
	}

	return r, nil
}

// Change the position in the file to the specified offset. If whence if 0, the
// offset is relative to the beggining of the file; if whence is 1, it's
// relative to the current offset; if whence is 2 it's relative to the end of
// the file. Any other value will result in an error. Returns the new offset
// and an error, if any.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.isdir() {
		return 0, os.EISDIR
	}

	newoff := offset
	switch whence {
	case 0:
		break
	case 1:
		cp, err := f.Tell()
		if err != nil {
			return newoff + cp, err
		}
		newoff += cp
	case 2:
		eof, err := f.Length()
		if err != nil {
			return eof + newoff, err
		}
		newoff += eof
	default:
		return newoff, errors.New(fmt.Sprintf("Unknown value for whence: %v", whence))
	}

	r := int64(C.PHYSFS_seek(f.cfile, C.PHYSFS_uint64(newoff)))

	if r == 0 {
		return newoff, errors.New(GetLastError())
	}

	return newoff, nil
}

// Returns the total length of the file and an error, if any.
func (f *File) Length() (int64, error) {
	if f.isdir() {
		return 0, os.EISDIR
	}

	r := int64(C.PHYSFS_fileLength(f.cfile))

	if r == -1 {
		return r, errors.New(GetLastError())
	}

	return r, nil
}

// Set up buffering for a PhysicsFS file handle. The following is copied almost
// verbatim from the PhysicsFS API reference:
//
// Define an i/o buffer for a file handle. A memory block of size bytes
// will be allocated and associated with the file. For files opened for reading,
// up to size bytes are read from the file and stored in the internal
// buffer. Calls to File.Read() will pull from this buffer until it is empty,
// and then refill it for more reading. Note that compressed files, like ZIP
// archives, will decompress while buffering, so this can be handy for
// offsetting CPU-intensive operations. The buffer isn't filled until you do
// your next read. For files opened for writing, data will be buffered to memory
// until the buffer is full or the buffer is flushed. Closing a handle
// implicitly causes a flush...check your return values! Seeking, etc
// transparently accounts for buffering. You can resize an existing buffer by
// calling this function more than once on the same file. Setting the buffer
// size to zero will free an existing buffer. PhysicsFS file handles are
// unbuffered by default. Please check the return value of this function!
// Failures can include not being able to seek backwards in a read-only file
// when removing the buffer, not being able to allocate the buffer, and not
// being able to flush the buffer to disk, among other unexpected problems.
func (f *File) SetBuffer(size uint64) error {
	if f.isdir() {
		return os.EISDIR
	}

	if int(C.PHYSFS_setBuffer(f.cfile, C.PHYSFS_uint64(size))) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// Flush the buffer of a buffered file. If the file was only opened for reading
// or is unbuffered this will do nothing successfully. Returns an error, if any.
func (f *File) Flush() error {
	if f.isdir() {
		return os.EISDIR
	}

	if int(C.PHYSFS_flush(f.cfile)) != 0 {
		return nil
	}

	return errors.New(GetLastError())
}

// A synonym for File.Flush(). Exactly the same.
func (f *File) Sync() error {
	return f.Flush()
}

// TODO: Make File.Stat() and File.Readdir() actually work correctly.

func (f *File) Stat() (fi os.FileInfo, err error) {
	size, err := f.Length()
	if err != nil {
		return
	}

	return fileInfo{
		f.name,
		size,
	}, nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	if !f.isdir() {
		return nil, os.ENOTDIR
	}
	if f.read < 0 {
		return nil, os.EINVAL
	}

	files, err := EnumerateFiles(f.name)
	if err != nil {
		return nil, err
	}

	if count < 0 {
		count = len(files)
	}
	if len(files)-f.read < count {
		count = len(files) - f.read
	}

	fi := make([]os.FileInfo, 0, count)
	for i := range files[len(files)-count:] {
		file, err := Open(f.name + "/" + files[i])
		if err != nil {
			return nil, err
		}
		info, err := file.Stat()
		if err != nil {
			return nil, err
		}
		fi = append(fi, info)
		file.Close()
	}

	f.read += count

	return fi, nil
}

type fileInfo struct {
	name string
	size int64
}

func (fi fileInfo) Name() string {
	return fi.name
}

func (fi fileInfo) Size() int64 {
	return fi.size
}

func (fi fileInfo) Mode() os.FileMode {
	if fi.IsDir() {
		return os.ModeDir
	}

	return 0644
}

func (fi fileInfo) ModTime() time.Time {
	mt, err := GetLastModTime(fi.name)
	if err != nil {
		return time.Time{}
	}

	return mt
}

func (fi fileInfo) IsDir() bool {
	return IsDirectory(fi.name)
}

type fileSystem struct{}

// Returns a simple implementation of http.FileSystem that simply opens the
// specified PhysicsFS file.
// Currently pointless due to the fact the File doesn't satisfy http.File.
func FileSystem() http.FileSystem {
	return new(fileSystem)
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	return Open(name)
}
