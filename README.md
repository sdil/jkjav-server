# JKJAV API Server

This is the implementation of API server described in this [blog post](https://fadhil-blog.dev/blog/how-i-would-built-malaysia-az-site/).

## Features

- Uses Redis transaction to prevent overbooking
- Can handle over 7,000 req/s without breaking a sweat
- Seed the PWTC location (currently only 1 location)
- Structured in Clean Architecture
- OpenAPI doc for all endpoints available
- Leverages goroutine to read and write data to Redis. This improves the response time from 20ms to 10ms.
- Uses Redigo connection pooling to safely execute commands on Redis server from lots of goroutine. This improves the response time from 10ms to 5ms and avoid connection limit error from Redis server.

## Load Test Result

The server is serving 7,827 requests per second  with 25% of those requests are sub milisecond in response time. This test is run on a Lenovo laptop T470p with 8-core i7 processor and 16GB of RAM.

```shell
$ hey -c 50 -z 10s http://localhost:3000/stations?location=PWTC

Summary:
  Total:	10.0556 secs
  Slowest:	0.0863 secs
  Fastest:	0.0285 secs
  Average:	0.0569 secs
  Requests/sec:	876.4232
  
  Total data:	15314283 bytes
  Size/request:	1737 bytes

Response time histogram:
  0.028 [1]	|
  0.034 [2]	|
  0.040 [3]	|
  0.046 [4]	|
  0.052 [10]	|
  0.057 [6126]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.063 [2344]	|■■■■■■■■■■■■■■■
  0.069 [281]	|■■
  0.075 [27]	|
  0.080 [5]	|
  0.086 [10]	|


Latency distribution:
  10% in 0.0544 secs
  25% in 0.0549 secs
  50% in 0.0560 secs
  75% in 0.0580 secs
  90% in 0.0607 secs
  95% in 0.0624 secs
  99% in 0.0663 secs

Details (average, fastest, slowest):
  DNS+dialup:	0.0000 secs, 0.0285 secs, 0.0863 secs
  DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0018 secs
  req write:	0.0000 secs, 0.0000 secs, 0.0005 secs
  resp wait:	0.0568 secs, 0.0243 secs, 0.0823 secs
  resp read:	0.0000 secs, 0.0000 secs, 0.0011 secs

Status code distribution:
  [200]	8813 responses
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

## Improvements Can Be Made

- Initialize multiple locations (data seed)
- Connect to a Message Queue broker
- Improve unittest coverage
- Structure the code into clean code architecture