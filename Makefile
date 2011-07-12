MAKE=gomake

SRCDIR:=physfs

.PHONY: all install test clean fmt

all:
	$(MAKE) -C $(SRCDIR) $@

install install.clean: all
	$(MAKE) -C $(SRCDIR) $@

test clean:
	$(MAKE) -C $(SRCDIR) $@

fmt:
	$(MAKE) -C $(SRCDIR) $@
