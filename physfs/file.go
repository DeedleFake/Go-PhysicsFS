package physfs

import(
	"os"
	"fmt"
	"unsafe"
)

// #include <stdlib.h>
// #include <physfs.h>
import "C"

// A handle for PhysicsFS. Designed to be as close to compatible as possible
// with os.File.
type File C.PHYSFS_File

// Open the named file with the mode specified by flag. The only supported
// arguments to flag 'os.O_RDONLY', 'os.O_WRONLY', and 'os.O_WRONLY
// | os.APPEND'. Returns the file and an error, if any.
func Open(name string, flag int) (f *File, err os.Error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	switch flag {
		case os.O_RDONLY:
			f = (*File)(C.PHYSFS_openRead(cname))
		case os.O_WRONLY:
			f = (*File)(C.PHYSFS_openWrite(cname))
		case os.O_WRONLY | os.O_APPEND:
			f = (*File)(C.PHYSFS_openAppend(cname))
		default:
			return nil, os.NewError("Unknown flag(s).")
	}

	if f == nil {
		return f, os.NewError(GetLastError())
	}

	return f, nil
}

// Close the file, release related resources. Returns an error, if any.
func (f *File)Close() (os.Error) {
	if int(C.PHYSFS_close((*C.PHYSFS_File)(f))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

// Read up to len(buf) bytes from the file into buf. Returns the number of bytes
// read and an error, if any.
func (f *File)Read(buf []byte) (n int, err os.Error) {
	n = int(C.PHYSFS_read((*C.PHYSFS_File)(f), unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))

	if n == -1 {
		err = os.NewError(GetLastError())
	}

	if f.EOF() {
		err = os.EOF
	}

	return n, err
}

// Write the bytes in buf to the file. Returns the number of bytes written and
// an error, if any.
func (f *File)Write(buf []byte) (n int, err os.Error) {
	n = int(C.PHYSFS_write((*C.PHYSFS_File)(f), unsafe.Pointer(&buf[0]), 1, C.PHYSFS_uint32(len(buf))))

	if n == -1 {
		return n, os.NewError(GetLastError())
	}

	return n, nil
}

// Returns a boolean indicating whether or not the end of the file has been
// reached.
func (f *File)EOF() (bool) {
	if int(C.PHYSFS_eof((*C.PHYSFS_File)(f))) != 0 {
		return true
	}

	return false
}

// Returns a number indication the current position in the file, and an error,
// if any.
func (f *File)Tell() (int64, os.Error) {
	r := int64(C.PHYSFS_tell((*C.PHYSFS_File)(f)))
	if r == -1 {
		return r, os.NewError(GetLastError())
	}

	return r, nil
}

// Change the position in the file to the specified offset. If whence if 0, the
// offset is relative to the beggining of the file; if whence is 1, it's
// relative to the current offset; if whence is 2 it's relative to the end of
// the file. Any other value will result in an error. Returns the new offset
// and an error, if any.
func (f *File)Seek(offset int64, whence int) (int64, os.Error) {
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
			return newoff, os.NewError(fmt.Sprintf("Unknown value for whence: %v", whence))
	}

	r := int64(C.PHYSFS_seek((*C.PHYSFS_File)(f), C.PHYSFS_uint64(newoff)))

	if r == 0 {
		return newoff, os.NewError(GetLastError())
	}

	return newoff, nil
}

// Returns the total length of the file and an error, if any.
func (f *File)Length() (int64, os.Error) {
	r := int64(C.PHYSFS_fileLength((*C.PHYSFS_File)(f)))

	if r == -1 {
		return r, os.NewError(GetLastError())
	}

	return r, nil
}

// Set up buffering for a PhysicsFS file handle. The following is copied almost
// verbatim from the PhysicsFS API reference:
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
// size to zero will free an existing buffer. PhysicsFS file handles are unbuffered by default. Please check the return value of this function! Failures can
// include not being able to seek backwards in a read-only file when removing
// the buffer, not being able to allocate the buffer, and not being able to
// flush the buffer to disk, among other unexpected problems.
func (f *File)SetBuffer(size uint64) (os.Error) {
	if int(C.PHYSFS_setBuffer((*C.PHYSFS_File)(f), C.PHYSFS_uint64(size))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

// Flush the buffer of a buffered file. If the file was only opened for reading
// or is unbuffered this will do nothing successfully. Returns an error, if any.
func (f *File)Flush() (os.Error) {
	if int(C.PHYSFS_flush((*C.PHYSFS_File)(f))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

// A synonym for File.Flush(). Exactly the same.
func (f *File)Sync() (os.Error) {
	return f.Flush()
}
