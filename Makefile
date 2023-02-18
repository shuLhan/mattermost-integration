## Copyright 2023 M. Sulhan <ms@kilabit.info>. All rights reserved.
## Use of this source code is governed by a BSD-style
## license that can be found in the LICENSE file.

COVER_OUT:=cover.out
COVER_HTML:=cover.html

.PHONY: all test lint clean

all: test lint

test: CGO_ENABLED=1
test:
	go test -race -coverprofile=$(COVER_OUT) ./...
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

lint:
	-golangci-lint run ./...
	-fieldalignment ./...

clean:
	rm -f $(COVER_OUT) $(COVER_HTML)
