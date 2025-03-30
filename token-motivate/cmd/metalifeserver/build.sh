#!/bin/bash

source ../../env.sh
#export CGO_ENABLED=0

echo -e "\e[1;31m VERSION \e[0m" \\t $VERSION
echo -e "\e[1;31m BUILD_DATE \e[0m" \\t $BUILD_DATE
echo -e "\e[1;31m GIT_COMMIT \e[0m" \\t $GIT_COMMIT
echo -e "\e[1;31m GO_VERSION \e[0m" \\t $GO_VERSION

go build -ldflags "-w -X main.Version=$VERSION -X main.Build=$BUILD_DATE -X main.GitCommit=$GIT_COMMIT -X main.GoVersion=$GO_VERSION "

cp metalifeserver $GOPATH/bin