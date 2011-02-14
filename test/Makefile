include $(GOROOT)/src/Make.inc

TARG=test
GOFILES=$(wildcard *.go)

CLEANFILES+=file2

GC=$Og -I ../_obj
LD=$Ol -L ../_obj

include $(GOROOT)/src/Make.cmd
