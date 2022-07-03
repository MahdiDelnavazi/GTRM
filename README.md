# GTRM 

GTRM using Golang-Transaction-Rabbitmq-Mongodb to handle transaction and ACID in mongodb

## Problem
when we want to change one field (in this example is a count of the bench document) we shuld use transaction to dont get race condition and bad data integrity. for handle this we use rabbitmq to queue all process and with this we dont get any bad data intigrity.

## Setting up the environment

For start server :
```
$ make server
```

For start consumer (rabbitmq) :
```
$ make ampqServer
```

## Benchmark

appache benchmark :
```
ab -k -u json.txt -T aplication/json -n 1000000 -c 32 -t 20  http://127.0.0.1:8080/mongo-increment
```

wrk :
```
wrk -t8 -c10000 -d20s -s script.lua --latency http://127.0.0.1:8080/mongo-increment
```
