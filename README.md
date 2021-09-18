# JKJAV API Server

This is the implementation of API server described in this [blog post](https://fadhil-blog.dev/blog/how-i-would-built-malaysia-az-site/).

## Features

- Uses Redis transaction to prevent overbooking and ensure atomicity
- Can handle over 7,000 req/s without breaking a sweat
- Seed the vaccine station (currently only 1 location)
- Structured in Clean Architecture
- OpenAPI doc for all endpoints available
- Leverages goroutine to read and write data to Redis. This improves the response time from 20ms to 10ms, 300 req/s to 800 req/s.
- Uses Redigo connection pooling to safely execute commands on Redis server from lots of goroutine. This improves the response time from 10ms to 5ms, 800 req/s to 1200 req/s and avoid connection limit error from Redis server.
- Connects to Kafka-compatible broker
- Comes with docker-compose yaml file to ease environment creation

## Load Test Result

The server is serving 7,827 requests per second  with 25% of those requests are sub milisecond in response time. This test is run on a Lenovo laptop T470p with 8-core i7 processor and 16GB of RAM.

```shell
$ hey -c 50 -z 10s http://localhost:3000/stations?location=PWTC

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
$hey -c 10 -z 1s http://localhost:3000/submit

Summary:
  Total:	1.0017 secs
  Slowest:	0.0082 secs
  Fastest:	0.0006 secs
  Average:	0.0013 secs
  Requests/sec:	7827.9083
  
  Total data:	740366 bytes
  Size/request:	94 bytes

Response time histogram:
  0.001 [1]	|
  0.001 [6394]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.002 [1070]	|■■■■■■■
  0.003 [101]	|■
  0.004 [20]	|
  0.004 [38]	|
  0.005 [69]	|
  0.006 [73]	|
  0.007 [34]	|
  0.007 [24]	|
  0.008 [17]	|


Latency distribution:
  10% in 0.0008 secs
  25% in 0.0009 secs
  50% in 0.0011 secs
  75% in 0.0013 secs
  90% in 0.0016 secs
  95% in 0.0021 secs
  99% in 0.0059 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0006 secs, 0.0082 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0003 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0004 secs
  resp wait:	0.0012 secs, 0.0005 secs, 0.0082 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0014 secs

Status code distribution:
  [200]	4292 responses
  [418]	3549 responses
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

## Caveat
- get station endpoint will give the result not in order because it's done by goroutines asynchronously

## Improvements Can Be Made

- Initialize multiple locations (data seed)
- Connect to a Message Queue broker
- Improve unittest coverage