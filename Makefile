PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

.PHONY: build install uninstall man test

build:
	go build -o adl ./cmd/adl

install: build
	install -d $(BINDIR) $(MANDIR)
	install -m 755 adl $(BINDIR)/adl
	install -m 644 man/adl.1 $(MANDIR)/adl.1

uninstall:
	rm -f $(BINDIR)/adl $(MANDIR)/adl.1

man:
	man -l man/adl.1

test:
	go test ./...
