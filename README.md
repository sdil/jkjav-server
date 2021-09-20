# JKJAV API Server

This is the implementation of API server described in this [blog post](https://fadhil-blog.dev/blog/how-i-would-built-malaysia-az-site/).

## Features

- Uses Redis transaction to prevent overbooking and ensure atomicity
- Seed the vaccine station (currently only 1 location)
- Structured in Clean Architecture
- OpenAPI doc for all endpoints available
- Leverages goroutine to read and write data to Redis. This improves the response time from 20ms to 10ms, 300 req/s to 800 req/s.
- Uses Redigo connection pooling to safely execute commands on Redis server from lots of goroutines in parallel. This improves the response time from 10ms to 5ms, 800 req/s to 1200 req/s and avoid connection limit error from Redis server.
- Connects to Redpanda, a Kafka-compatible broker and publishes message after a new booking is made
- Comes with docker-compose yaml file to ease environment creation
- Unit test built with Mockery & table driven test

## Load Test Result

The server is serving over 8,600 requests per second with 95% of those requests are <15ms response time. This test is run on a Lenovo laptop T470p with 8-core i7 processor and 16GB of RAM.

```shell
$ hey -c 50 -z 10s http://localhost:3000/stations/PWTC

Summary:
  Total:	10.0385 secs
  Slowest:	0.1060 secs
  Fastest:	0.0079 secs
  Average:	0.0403 secs
  Requests/sec:	1236.0401
  
  Total data:	21488345 bytes
  Size/request:	1731 bytes

Response time histogram:
  0.008 [1]	|
  0.018 [6]	|
  0.028 [10]	|
  0.037 [179]	|■
  0.047 [12088]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.057 [47]	|
  0.067 [20]	|
  0.077 [7]	|
  0.086 [7]	|
  0.096 [27]	|
  0.106 [16]	|


Latency distribution:
  10% in 0.0384 secs
  25% in 0.0391 secs
  50% in 0.0399 secs
  75% in 0.0410 secs
  90% in 0.0420 secs
  95% in 0.0427 secs
  99% in 0.0472 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0079 secs, 0.1060 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0031 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0011 secs
  resp wait:	0.0403 secs, 0.0048 secs, 0.1060 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0012 secs

Status code distribution:
  [200]	12408 responses
```


```shell
$ hey -c 100 -z 5s -T "application/json" -m POST -d '{"address":"Kuala Lumpur","date":"20210516","firstName":"Fadhil","lastName":"Yaacob","location":"PWTC","mysejahteraId":"900127015527","phoneNumber":"0123456789"}' http://localhost:3000/booking

Summary:
  Total:	5.0093 secs
  Slowest:	0.1446 secs
  Fastest:	0.0008 secs
  Average:	0.0116 secs
  Requests/sec:	8609.4523
  
  Total data:	2245918 bytes
  Size/request:	52 bytes

Response time histogram:
  0.001 [1]	|
  0.015 [41010]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.030 [903]	|■
  0.044 [112]	|
  0.058 [17]	|
  0.073 [117]	|
  0.087 [193]	|
  0.101 [373]	|
  0.116 [201]	|
  0.130 [100]	|
  0.145 [100]	|


Latency distribution:
  10% in 0.0073 secs
  25% in 0.0080 secs
  50% in 0.0089 secs
  75% in 0.0103 secs
  90% in 0.0121 secs
  95% in 0.0150 secs
  99% in 0.1000 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0008 secs, 0.1446 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0023 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0034 secs
  resp wait:	0.0115 secs, 0.0008 secs, 0.1445 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0071 secs

Status code distribution:
  [200]	7268 responses
  [400]	35859 responses

```

## How to run

```shell
make run

# In another terminal
make load-test
make redis-cli
```

Run this in redis-cli
```
keys user:*
keys location:*
```

## Caveats
- get station endpoint will give the result not ordered in date because the data fetching are done by goroutines in parallel

## Improvements Can Be Made

- Initialize multiple locations (data seed)
- Improve unittest coverage
- Add telemetry
