MAKE=gomake

SRCDIR:=physfs

.PHONY: all install clean install.clean test fmt

all test clean fmt:
	$(MAKE) -C $(SRCDIR) $@

install install.clean: all
	$(MAKE) -C $(SRCDIR) $@
