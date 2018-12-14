#! /bin/sh
REDIS=`netstat -an | grep 6379`
if [ -z "$REDIS" ]; then 
    echo "redis is not strated!"
    docker run -d -p 6379:6379 --name="redisdb" redis 
fi
go run dataServer/dataServer.go