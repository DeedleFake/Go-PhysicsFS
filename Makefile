include $(GOROOT)/src/Make.inc

TARG=physfs
CGOFILES=\
		physfs.go\
		file.go\

CGO_LDFLAGS=-lphysfs

include $(GOROOT)/src/Make.pkg
