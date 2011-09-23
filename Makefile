MAKE=gomake

SRCDIR:=physfs

.PHONY: all install clean install.clean test fmt

all test clean:
	$(MAKE) -C $(SRCDIR) $@

install install.clean: all
	$(MAKE) -C $(SRCDIR) $@

fmt:
	$(MAKE) -C $(SRCDIR) $@
	$(MAKE) -C test $@
	$(MAKE) -C httptest $@
