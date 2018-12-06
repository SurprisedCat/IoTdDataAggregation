#! /bin/sh
REDIS=`docker ps  | grep redis | awk {'print $2'}`
if [ -z "$REDIS" ]; then 
    echo "redis is not strated!"
    docker run -d -p 6379:6379 --name="redisdb" redis 
fi
go run dataServer/dataServer.go