#!/bin/bash

case "$1" in
  start)
    pushd /Users/ahmedshawky/learn/distributed-tracing
    (cd platform-demo; rails s; echo $! > pid1)
    pushd /Users/ahmedshawky/learn/distributed-tracing
    (cd vendor-availability-service; rails s  -p 3001; echo $! > pid2)
    pushd /Users/ahmedshawky/learn/distributed-tracing
    (cd payment-service; rails s  -p 3002; echo $! > pid3)
    pushd /Users/ahmedshawky/learn/distributed-tracing
    (cd order-transmission-servcie; rails s  -p 3003; echo $! > pid4)
    pushd /Users/ahmedshawky/learn/distributed-tracing
    (cd time-estimation-service; rails s -p 3004; echo $! > pid5)

    popd
    ;;

  stop)
    kill $(cat pid1)
    kill $(cat pid2)
    kill $(cat pid3)
    kill $(cat pid4)
    kill $(cat pid5)
    rm pid1 pid2 pid3 pid4 pid5
    ;;

  *)
    echo "Usage: $0 {start|stop}"
    exit 1
    ;;
esac

exit 0
