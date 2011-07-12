include $(GOROOT)/src/Make.inc

TARG=httptest
GOFILES=$(wildcard *.go)

CLEANFILES+=\
			file2\
			test-dir\

GC=$Og -I ../physfs/_obj
LD=$Ol -L ../physfs/_obj

include $(GOROOT)/src/Make.cmd

.PHONY: fmt

fmt:
	gofmt -w $(GOFILES)
