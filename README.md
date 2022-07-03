# GTRM

GTRM using :
* [**Golang**](https://go.dev/)
* [**RabbitMQ**](https://www.rabbitmq.com/)
* [**MongoDB**](https://www.mongodb.com/)
* [**Gin Gonic**](https://github.com/gin-gonic/gin)

## Motivation
### ACID : 
Collections of operations that form a single logical unit of work are called transactions and the database system must ensure proper execution of transactions and the ACID database transaction model ensures that a performed transaction is always consistent.
### BASE : 
The rise in popularity of NoSQL databases provided a flexible and fluidity with ease to manipulate data and as a result, a new database model was designed, reflecting these properties. The acronym BASE is slightly more confusing than ACID but however, the words behind it suggest ways in which the BASE model is different and acronym BASE stands for: 1 - Basically Available 2 - soft state 3 - Eventually Consistent



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
