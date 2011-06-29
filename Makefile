MAKE=gomake

SRCDIR:=physfs

.PHONY: all install test clean

all:
	$(MAKE) -C $(SRCDIR) $@

install install.clean: all
	$(MAKE) -C $(SRCDIR) $@

test clean:
	$(MAKE) -C $(SRCDIR) $@
