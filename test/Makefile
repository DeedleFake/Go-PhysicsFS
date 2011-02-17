include $(GOROOT)/src/Make.inc

TARG=test
GOFILES=$(wildcard *.go)

CLEANFILES+=\
			file2\
			test-dir\

GC=$Og -I ../physfs/_obj
LD=$Ol -L ../physfs/_obj

.PHONY: all

all: test-dir $(TARG)

test-dir:
	mkdir -p test-dir
	touch test-dir/{this,is,a,test}

include $(GOROOT)/src/Make.cmd
