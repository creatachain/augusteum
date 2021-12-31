#! /bin/bash
set -e

# These tests spawn the counter app and server by execing the MSM_APP command and run some simple client tests against it

# Get the directory of where this script is.
export PATH="$GOBIN:$PATH"
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# Change into that dir because we expect that.
cd "$DIR"

echo "RUN COUNTER OVER SOCKET"
# test golang counter
MSM_APP="counter" go run -mod=readonly ./*.go
echo "----------------------"


echo "RUN COUNTER OVER GRPC"
# test golang counter via grpc
MSM_APP="counter --msm=grpc" MSM="grpc" go run -mod=readonly ./*.go
echo "----------------------"

# test nodejs counter
# TODO: fix node app
#MSM_APP="node $GOPATH/src/github.com/creatachain/js-msm/example/app.js" go test -test.run TestCounter
