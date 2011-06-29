include $(GOROOT)/src/Make.inc

TARG=physfs
CGOFILES=\
         physfs.go\
         file.go\

CLEANFILES+=\
			doc.html\

CGO_OFILES+=\
			wrapcb.o

include $(GOROOT)/src/Make.pkg

doc:
	godoc -html -path=. . > doc.html
