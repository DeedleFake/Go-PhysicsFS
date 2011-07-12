MAKE=gomake

SRCDIR:=physfs

.PHONY: all install test clean fmt

all test clean fmt:
	$(MAKE) -C $(SRCDIR) $@

install install.clean: all
	$(MAKE) -C $(SRCDIR) $@
