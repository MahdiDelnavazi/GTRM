# GTRM 

GTRM using Golang-Transaction-RabbitMQ-MongoDB to handle transaction and ACID in MongoDB

## Problem
When we want to change one field (in this example is a count of the bench document) we should use transaction to don't get race condition and bad data integrity. For handle this we use RabbitMQ to queue all process and with this we don't get any bad data integrity.

## Setting up the environment

For start server :
```
$ make server
```

For start consumer (RabbitMQ) :
```
$ make ampqServer
```

## Benchmark

Appache benchmark :
```
ab -k -u json.txt -T aplication/json -n 1000000 -c 32 -t 20  http://127.0.0.1:8080/mongo-increment
```

wrk :
```
wrk -t8 -c10000 -d20s -s script.lua --latency http://127.0.0.1:8080/mongo-increment
```
