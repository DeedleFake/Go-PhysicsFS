MAKE=gomake

SRCDIR:=physfs

.PHONY: all install test clean

all:
	$(MAKE) -C $(SRCDIR) $@

install: all
	$(MAKE) -C $(SRCDIR) $@

test:
	$(MAKE) -C $(SRCDIR) $@

clean:
	$(MAKE) -C $(SRCDIR) $@
