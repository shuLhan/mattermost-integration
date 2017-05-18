GOPATH_SHARED_DEPS:= $(GOPATH)/src/github.com/Sirupsen/logrus \

COVER_OUT:=	cover.out
COVER_HTML:=	cover.html

SRC:=		$(shell go list -f '{{join .GoFiles " "}}')
TEST_SRC:=	$(shell go list -f '{{join .TestGoFiles " "}}')

TARGET:=	logrus.a
TARGET_PKG:=	$(shell go list -f '$(GOPATH)/pkg/{{context.GOOS}}_{{context.GOARCH}}/{{.ImportPath}}.a')

.PHONY: all
.PHONY: clean distclean
.PHONY: check check-all
.PHONY: test coverbrowse
.PHONY: install

##
## MAIN TARGETS
##

all: install

##
## DEPENDENCIES
##

$(GOPATH)/src/github.com/Sirupsen/logrus:
	go get -u github.com/Sirupsen/logrus

deps: $(GOPATH_SHARED_DEPS)

##
## CLEAN
##

clean:
	rm -f $(COVER_OUT) $(COVER_HTML)

distclean: clean

##
## LINT
##

check: $(SRC) $(TEST_SRC)
	gometalinter --concurrency=2 --vendor --deadline=240s --disable=gotype ./...

check-all: $(SRC) $(TEST_SRC)
	gometalinter --concurrency=2 --vendor --deadline=240s --disable=gotype --enable-all ./...

check-errors: $(SRC)
	gometalinter --concurrency=2 --vendor --deadline=240s --disable=gotype --fast --errors  ./...

##
## TEST
##

$(COVER_OUT): $(SRC) $(TEST_SRC)
	go test -v -coverprofile=$@

$(COVER_HTML): $(COVER_OUT)
	go tool cover -html=$(COVER_OUT) -o $(COVER_HTML)

test: deps check-errors $(COVER_HTML)

coverbrowse: $(COVER_HTML)
	xdg-open $(COVER_HTML) &

##
## BUILD
##

$(TARGET): test $(SRC) $(TEST_SRC)
	go build -o $@

##
## INSTALL and DEPLOY
##

$(TARGET_PKG): test $(SRC) $(TEST_SRC)
	@echo ">>> Target pkg: $@"
	go install -a

install: $(TARGET_PKG)
