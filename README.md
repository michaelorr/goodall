        ____                 _       _ _
       / ___| ___   ___   __| | __ _| | |
      | |  _ / _ \ / _ \ / _` |/ _` | | |
      | |_| | (_) | (_) | (_| | (_| | | |
       \____|\___/ \___/ \__,_|\__,_|_|_|


# Metrics gathering and reporting

## What is this?

TODO

## Naming things is hard

The name is derived from Jane Goodall, famed primatologist and anthropologist
who became famous for her 55 year observational study of chimpanzees in
Tanzania. Also, there is a trend of incorporting the letters `Go` in Go based
projects ;-)

> https://www.wikiwand.com/en/Jane_Goodall

## Building

Run `make` or `make build` to compile your app.  This will use a Docker image
to build your app, with the current directory volume-mounted into place.  This
will store incremental state for the fastest possible build.  Run `make
all-build` to build for all architectures.

Run `make clean` to clean up.

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

## Adding new Metrics

To add a new metric, add an entry to BucketMap in `pkg/metrics/metrics.go`.
This map is of the form `"key": value` where `key` is a string and `value` is a
function with the signature `func(string, chan *DataPoint, chan error)`. The
function should send along the channel a pointer to an instance of
metrics.Datapoint which contains the BucketName that was passed in and the
result value that should be stored in the db. This method will be called every
metrics.Interval which is an instance of time.Duration.

Goodall depends on https://github.com/shirou/gopsutil for gathering system
metrics. This library can read environment variable if the location of `/proc`
`/etc` or `/sys` are different for your target system. Check out
https://godoc.org/github.com/shirou/gopsutil for more documentation or other
metrics that can be easily added.

## Notes

* Only the amd64/darwin and amd64/linux have been tested.
* Other OSs and ARCHs may build successfully but have not been tested and
likely do not work fully due to lack of gopsutil support on Windows.

## Error Response Codes

### 0 - 9

* 1: There was an error opening the database file
* 2: There was a problem initializing the database buckets

### 10-19

* TODO

## CLI Params

* TODO
* Interval
* DB filename
* Data Cleanup

## Resource Utilization

If this tool is intended to monitor resources, it shouldn't be resource heavy
itself. I've measured steady state usage with default configuration (gather
metrics every 1s, store for 4h). Disk utilization by the database is roughly
25M with steady state usage of 6-7M resident memory while CPU utilization
remains reasonably low. If the metric interval is increased to 1ms, the disk
utilization will naturally increase dramatically as will CPU utilization, but
overall resident memory usage will remain quite stable.

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
