export MAKEFILE=MakeFile
# get base folder name
FOLDERNAME=$(shell basename "$(PWD)")
# This how we want to name the binary output
BINARY=dist/favagateway
# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
# git branch name
BRANCHANME=`git rev-parse --abbrev-ref HEAD`
VERSION=`git describe --tags`
BUILD=`date +'%Y-%m-%dT%H:%M:%S'`
BUILDSAVE=`date +'%Y_%m_%dT%H_%M'`
# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Branch=${BRANCHANME}"
# Default target that depends on clean and build
all: clean build
# remove build file if exists
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	if [ -f ${BINARY}_MyOS ] ; then rm ${BINARY}_MyOS ; fi
# Builds the project
build:
	GOOS=darwin GOARCH=amd64 go build -o build ${LDFLAGS} -o ${BINARY}_${BUILDSAVE}_mac
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY}_${BUILDSAVE}_linux
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY}_${BUILDSAVE}_win.exe