include $(GOROOT)/src/Make.inc

TARG=physfs
CGOFILES=\
         physfs.go\
         file.go\

CGO_LDFLAGS=-lphysfs

CLEANFILES+=\
			doc.html\

include $(GOROOT)/src/Make.pkg

doc:
	godoc -html -path=. . > doc.html
