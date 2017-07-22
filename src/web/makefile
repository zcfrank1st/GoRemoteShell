PROJECT_DIR =  /Users/zcfrank1st/Desktop/git-svn/go-demo
#INSTALL_TARGET_DIR = /

export GOPATH := $(PROJECT_DIR)

all: build install clean

run:
	go run main.go

build:
	go build -i -o ../../bin/web/web main.go

install: build
#	TODO
#	cp ../../bin/web $(INSTALL_TARGET_DIR)/

clean:
	rm -rf ../../bin/ ../../pkg/
