        ____                 _       _ _
       / ___| ___   ___   __| | __ _| | |
      | |  _ / _ \ / _ \ / _` |/ _` | | |
      | |_| | (_) | (_) | (_| | (_| | | |
       \____|\___/ \___/ \__,_|\__,_|_|_|


# Metrics gathering and reporting

## Releases

[v0.1](https://github.com/michaelorr/goodall/releases/tag/v0.1-beta)

## What is this?

This is a metrics gathering and reporting system. The current implementation is
designed such that it can easily monitor one system. It will gather data on the
system and store it in a persistence layer provided by BoltDB
(github.com/boltdb/bolt) and will serve that data up in JSON format over an
http api. The api is described later in this document. The metrics that are
gathered by default are also described later in this document but this list
can easily be extended.
A future design goal would be to extract the metric gathering agent and data
reporting server into separate components so that they could be managed
independently and metrics could be gathered across a cluster of machines.

## Building

Run `make` or `make build` to compile your app.  This will use a Docker image
to build your app, with the current directory volume-mounted into place.  This
will store incremental state for the fastest possible build.  Run `make
all-build` to build for all architectures.

Run `make clean` to remove build artifacts.

## Testing

Run `make test` or `make test-verbose` to run the test suite. This will use a
Docker image to compile the test binary, install test dependencies, run the
test suite, run the `go fmt` linter and `go vet` static analyzer reporting the
results. Currently this make target only exercises the linux version of the
code.

Running `go test ./...` in the current directory should allow you to run
the tests on any OS/Architecture that is supported by Go, but you will need the
Go runtime and environment setup to do so. See https://golang.org/doc/install
for more information on how to get started.

NOTE: Sadly the test suite is currently non-existent. Running the make target
for the tests will still perform the `go fmt` and `go vet` checks but
unfortunately my recent efforts have been focused on finding the memory leak.

## API

* `/latest`: This end point will return a payload which contains one data entry
per metric gathered. Each metric will be the most recent of its type.
* There is also an endpoint for each named metric, i.e. `/disk_used`,
`/mem_total`, etc. These endpoints return the 10 most recent entries for this
metric type.

## JSON Body Spec

The JSON payload will consist of the following structure:

```
{
    "Timestamp": "2016-10-29T05:04:16.395",
    "Metrics": [
        {
            "Name": "system_load_15",
            "Value": 1.89,
            "Timestamp": "2016-10-29T05:04:11.513"
        },
        {
            "Name": "system_load_5",
            "Value": 1.76,
            "Timestamp": "2016-10-29T05:04:11.513"
        }
        ...
    ]
}
```

There is a timestamp for the overall payload body which represents the time at
which the payload was constructed and a series of metrics. Each metric contains
a key which indicates what data is in the metric, the actual value for the
metric, and the time at which that particular measurement was made.
Regardless of the API endpoint that was used to fetch the JSON payload, the
overall structure will be identical, but the contents will be tailored to suit
the request.

## Default Metrics

* `disk_used`
* `disk_free`
* `disk_total`
* `mem_used`
* `mem_available`
* `mem_total`
* `system_load_1`
* `system_load_15`
* `system_load_5`

The default metrics are provided by gopsutil
(https://godoc.org/github.com/shirou/gopsutil). More details can be found in
the documentation for that project.

## Adding new Metrics

To add a new metric, add an entry to BucketMap in `pkg/metrics/metrics.go`.
This map is of the form `"key": value` where `key` is a string and `value` is a
function with the signature `func(string, chan *DataPoint, chan error)`. The
function should send along the channel a pointer to an instance of
metrics.Datapoint which contains the Name if the bucket that was passed in and
the result value that should be stored in the db. This method will be called
every metrics.Interval which is an instance of time.Duration.

Goodall depends on https://github.com/shirou/gopsutil for gathering system
metrics. This library can read environment variable if the location of `/proc`
`/etc` or `/sys` are different for your target system. Check out
https://godoc.org/github.com/shirou/gopsutil for more documentation or other
metrics that can be easily added.

## Notes

* Only amd64/darwin and amd64/linux have been tested.
* Other OSs and ARCHs may build successfully but have not been tested and
likely do not work fully due to lack of gopsutil support on Windows.

## Error Response Codes

* `1`: There was an error opening the database file
* `2`: There was a problem initializing the database buckets

## Env Vars

* `GOODALL_PORT`: This is the port to serve data from. If unspecified, the
default value of `8080` is used. If the specified port is invalid, already in
use or if goodall fails to bind to the port, the server will fail to start.
* `GOODALL_COLLECTION_MS`: The collection interval expressed in milliseconds.
If unspecified, or unparseable by https://golang.org/pkg/strconv/#Atoi the
default value of `1000` is used.
* `GOODALL_RETENTION_MIN`: The retention time period expressed in minutes.
If unspecified, or unparseable by https://golang.org/pkg/strconv/#Atoi the
default value of `40` is used.
* `GOODALL_DB_PATH`: This is the string filepath to the location of the db
BoltDB file.
If unspecified, the default of `goodall.db` is used. If the DB file does not
exist when the service starts, the db file will be created.

NOTE: Goodall does not require the db file to exist but it does expect the
parent dir to exist and be writeable by the user. No validation is done on the
path and no fall-back is provided if the path or parent directory exists but is
un-writeable by the user or if the directory path does not exist.

## Resource Utilization

I've measured steady state usage with default configuration (gather metrics
every 1s, store for 4h). Disk utilization by the database is roughly 25MB while
CPU utilization remains reasonably low (fluctuations between 0% and momentary
blips of 5%). If the metric collection interval is increased, the disk
utilization will naturally increase as will CPU utilization. This will also
worsen the memory leak (see Known Issues). Increasing the retention period will
naturally have a linear relationship with disk usage by the DB.

## Benchmarks

While these benchmarks are purely anecodtal, they are an indicator of the
responsiveness that is possible with Goodall. I used ab
(https://linux.die.net/man/1/ab) which indicated that typical response times
were often in the neighborhood of 5-20ms with occasional outliers taking 3-10s
and typical concurrency on the order of 500-1500 req/sec (mean).

These measurements were taken while conducting various experiments of 5-100
concurrent requests and ranging from 5000 to 30000 total requests.

The following data was taken with 20 concurrent and 10000 total requests:
```
Concurrency Level:      20
Time taken for tests:   4.805 seconds
Complete requests:      10000
Failed requests:        0

Requests per second:    781.35 [#/sec] (mean)
Time per request:       9.609 [ms] (mean)

Percentage of the requests served within a certain time (ms)
    50%      4
    66%      4
    75%      5
    80%      5
    90%      6
    95%      8
    98%     13
    99%     20
   100%  13164 (longest request)
```

## Naming things is hard

The name is derived from Jane Goodall, famed primatologist and anthropologist
who became famous for her 55 year observational study of chimpanzees in
Tanzania. Also, there is a trend of incorporting the letters `Go` in Go based
projects ;-)

> `https://www.wikiwand.com/en/Jane_Goodall`

## Known Issues

There is a slow memory leak in the metric collection portion of code. I have
been unable to track it down precisely but it is made worse when the collection
period is increased. For this reason, it is not recommended to run Goodall with
a sub-second collection for any lengthy period of time. Running with the
default value of collecting metrics once per second for 1 hour, ram utilization
was ~10MB but after ~12 hours it was close to 100MB of ram. Not good.

## Wishlist

If I were going to operationalize this service as it is structured, there are
several things that are glossed over here. I would separate the gathering and
reporting responsibilties into separate binaries to allow for more granular
control over deployments. This would prevent the two processes from
communicating via the DB (realistically they shouldn't be doing this anyways),
this means that either a JSON or (more likely) a gRPC interface would be
established allowing the two processes to run on different hosts. This would
also allow this service to scale horizontally for metrics gathering across a
cluster or multiple reporters behind a load balancer to handle higher load.

There are also a lot of places where this service could hard-crash. That's not
ideal. I would want to wrap these places in some sort of retry logic,
potentially with an intelligent back-off that logs or alerts what is happening
or going wrong, but does not bring the service to it's knees.

Realistically, I would use a time series DB which is more geared to storing and
querying this type of data. InfluxDB is well suited to this problem space but
Graphite and Riak both support this. I chose BoltDB due to its simplicity in
not needing to run a separate process or daemon to achieve persistence.

Most importantly, if I were actually trying to build this functionality for
prod, I wouldn't build it at all. This is a well-travelled problem space and
there are many off the shelf solutions that are easily customizable and
extensible which provide much better reporting, graphing, and alerting than
any custom home-grown solution ever could. InfluxDB, graphite, statsd, datadog,
ELK, nagios, among a dozen others provide the same functionality as Goodall
plus much more.
