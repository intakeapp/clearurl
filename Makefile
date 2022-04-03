.PHONY: all build

all: build

build:
	go build ./cmd/clearurl

TMPDIR := $(shell mktemp -d)
fetch:
	git clone -b gh-pages git@github.com:ClearURLs/Rules.git $(TMPDIR)
	cd $(TMPDIR); git reset --hard 236f926b1eb833a47264779edb83db8d9497ff8d
	cp $(TMPDIR)/data.min.json .
