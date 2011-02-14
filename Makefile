include $(GOROOT)/src/Make.inc

TARG=physfs
CGOFILES=physfs.go
CGO_LDFLAGS=-lphysfs

include $(GOROOT)/src/Make.pkg
