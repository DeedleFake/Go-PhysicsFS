package physfs

import(
	"os"
	"unsafe"
)

// #include <physfs.h>
import "C"

type File C.PHYSFS_File

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

func (f *File)Length() (int64, os.Error) {
	r := int64(C.PHYSFS_fileLength((*C.PHYSFS_File)(f)))

	if r == -1 {
		return r, os.NewError(GetLastError())
	}

	return r, nil
}

func (f *File)SetBuffer(size uint64) (os.Error) {
	if int(C.PHYSFS_setBuffer((*C.PHYSFS_File)(f), C.PHYSFS_uint64(size))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func (f *File)Flush() (os.Error) {
	if int(C.PHYSFS_flush((*C.PHYSFS_File)(f))) != 0 {
		return nil
	}

	return os.NewError(GetLastError())
}

func (f *File)Sync() (os.Error) {
	return f.Flush()
}
