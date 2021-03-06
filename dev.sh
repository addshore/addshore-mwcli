#!/bin/bash

GOPATH=$(pwd)/vendor:$(pwd)
GOBIN=$(pwd)/bin

SCRIPTPATH=$(dirname $(realpath $0))
WORKINGPATH=$(echo $PWD)
cd $SCRIPTPATH

# Force remake of the mwdd files
rm internal/mwdd/files/files.go
make internal/mwdd/files/files.go

# Run from source from the origional directory
cd $WORKINGPATH
go run -ldflags "$(${SCRIPTPATH}/bin/govvv -flags)" ${SCRIPTPATH}/main.go $@