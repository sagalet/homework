# My homework

This is my homework. Design a server rate limitation system

##  Request
+ Golang > 1.11
+ Redis

## Server
Get the request from client and record the ip and the timestamp into redis
##### 1. Build & Run
+ go build -o out/server server/main.go
+ out/server -port=10034 -conf=server/redis.conf
<br /><br />
##### 2. Arguments
+ port : The port listened on server. The default is "10034".
+ conf : The redis information. The default is "redis.conf".


## Client
Send the request to the server and get the response
##### 1. Build & Run
+ go build -o out/client client/main.go
+ out/client -delay=2000 -conf=client/server.conf -count=120
<br /><br />
##### 2. Arguments
+ delay : The latency between each request . The default is 2000.
+ conf : The information of server. The default is "server.conf".
+ count : The number of reuqests sent. The default is 120

## Test
There are two tests
+ Call api to access redis server directly
+ Send http request to the server
##### 1. Run
+ go test -v test/common_test.go -port=10034 -redisip=127.0.0.1 -redisport=6379
<br /><br />
##### 2. Arguments
+ port : The port server listened. The default is 10034
+ redisip : The ip of redis server. The default is 127.0.0.1
+ redisport : The port of server listened. The default is 6379
