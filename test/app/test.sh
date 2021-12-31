#! /bin/bash
set -ex

#- kvstore over socket, curl
#- counter over socket, curl
#- counter over grpc, curl
#- counter over grpc, grpc

# TODO: install everything

export PATH="$GOBIN:$PATH"
export TMHOME=$HOME/.augusteum_app

function kvstore_over_socket(){
    rm -rf $TMHOME
    augusteum init
    echo "Starting kvstore_over_socket"
    msm-cli kvstore > /dev/null &
    pid_kvstore=$!
    augusteum node > augusteum.log &
    pid_augusteum=$!
    sleep 5

    echo "running test"
    bash test/app/kvstore_test.sh "KVStore over Socket"

    kill -9 $pid_kvstore $pid_augusteum
}

# start augusteum first
function kvstore_over_socket_reorder(){
    rm -rf $TMHOME
    augusteum init
    echo "Starting kvstore_over_socket_reorder (ie. start augusteum first)"
    augusteum node > augusteum.log &
    pid_augusteum=$!
    sleep 2
    msm-cli kvstore > /dev/null &
    pid_kvstore=$!
    sleep 5

    echo "running test"
    bash test/app/kvstore_test.sh "KVStore over Socket"

    kill -9 $pid_kvstore $pid_augusteum
}


function counter_over_socket() {
    rm -rf $TMHOME
    augusteum init
    echo "Starting counter_over_socket"
    msm-cli counter --serial > /dev/null &
    pid_counter=$!
    augusteum node > augusteum.log &
    pid_augusteum=$!
    sleep 5

    echo "running test"
    bash test/app/counter_test.sh "Counter over Socket"

    kill -9 $pid_counter $pid_augusteum
}

function counter_over_grpc() {
    rm -rf $TMHOME
    augusteum init
    echo "Starting counter_over_grpc"
    msm-cli counter --serial --msm grpc > /dev/null &
    pid_counter=$!
    augusteum node --msm grpc > augusteum.log &
    pid_augusteum=$!
    sleep 5

    echo "running test"
    bash test/app/counter_test.sh "Counter over GRPC"

    kill -9 $pid_counter $pid_augusteum
}

function counter_over_grpc_grpc() {
    rm -rf $TMHOME
    augusteum init
    echo "Starting counter_over_grpc_grpc (ie. with grpc broadcast_tx)"
    msm-cli counter --serial --msm grpc > /dev/null &
    pid_counter=$!
    sleep 1
    GRPC_PORT=36656
    augusteum node --msm grpc --rpc.grpc_laddr tcp://localhost:$GRPC_PORT > augusteum.log &
    pid_augusteum=$!
    sleep 5

    echo "running test"
    GRPC_BROADCAST_TX=true bash test/app/counter_test.sh "Counter over GRPC via GRPC BroadcastTx"

    kill -9 $pid_counter $pid_augusteum
}

case "$1" in 
    "kvstore_over_socket")
    kvstore_over_socket
    ;;
"kvstore_over_socket_reorder")
    kvstore_over_socket_reorder
    ;;
    "counter_over_socket")
    counter_over_socket
    ;;
"counter_over_grpc")
    counter_over_grpc
    ;;
    "counter_over_grpc_grpc")
    counter_over_grpc_grpc
    ;;
*)
    echo "Running all"
    kvstore_over_socket
    echo ""
    kvstore_over_socket_reorder
    echo ""
    counter_over_socket
    echo ""
    counter_over_grpc
    echo ""
    counter_over_grpc_grpc
esac

