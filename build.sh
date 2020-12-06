#!/bin/bash
 
#set -e
#set -x
 
PROJECT_NAME="qr-code service"
BINARY="qr-code"
CODE="main"

OUTPUT_DIR=output
GOOS=$(go env GOOS)
 
APP_NAME=${PROJECT_NAME}
#APP_VERSION=$(git log -1 --oneline)
#BUILD_VERSION=$(git log -1 --oneline)
let "BUILD_VERSION = $(< build_number.txt) + 1"
BUILD_TIME=$(date "+%F_%T.%N_%:::z")
APP_VERSION="0.1.${BUILD_VERSION}"
GIT_REVISION="0.0.1"
#GIT_REVISION=$(git rev-parse --short HEAD)
#GIT_BRANCH=$(git name-rev --name-only HEAD)
GIT_BRANCH=$(git branch | grep \*)
#GIT_BRANCH="master"
GO_VERSION=$(go version)


#-mod=vendor \
#CGO_ENABLED=0 go build -a -installsuffix cgo -v \
#GOARCH=amd64 GOOS=windows go build -ldflags "-s -X 'main.AppName=${APP_NAME}' \
go build -ldflags "-s -X 'main.AppName=${APP_NAME}' \
            -X 'main.AppVersion=${APP_VERSION}' \
            -X 'main.BuildVersion=${BUILD_VERSION}' \
            -X 'main.BuildTime=${BUILD_TIME}' \
            -X 'main.GitRevision=${GIT_REVISION}' \
            -X 'main.GitBranch=${GIT_BRANCH}' \
            -X 'main.GoVersion=${GO_VERSION}'" \
-o ${BINARY} ${CODE}.go

#if (( $? == 0 )); then

#GOARCH=amd64 GOOS=windows go build -ldflags "-s -X 'main.AppName=${APP_NAME}' \
#            -X 'main.AppVersion=${APP_VERSION}' \
#            -X 'main.BuildVersion=${BUILD_VERSION}' \
#            -X 'main.BuildTime=${BUILD_TIME}' \
#            -X 'main.GitRevision=${GIT_REVISION}' \
#            -X 'main.GitBranch=${GIT_BRANCH}' \
#            -X 'main.GoVersion=${GO_VERSION}'" \
#-o "${BINARY}.exe" ${CODE}.go
#else
#  $? = 1
#fi

if (( $? == 0 )); then
  echo "${BUILD_VERSION}" > build_number.txt
  echo "build app ${APP_VERSION} ok"
else
  echo "build error "
fi