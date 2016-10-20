# Metrics gathering and reporting

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

## Notes

* Only the amd64 arch for darwin and linux OSs have been fully tested.
* Other OSs and ARCHs build successfully but have not been tested.
* darwin/arm fails to build. I do not fully understand why at this time.

## Error Response Codes

### 0 - 9

* 1: There was an error opening the database file
* 2: There was a problem initializing the database buckets

### 10-19

* TODO

## Wishlist

If I were going to operationalize this service as it is structured, there are
several things that are glossed over here. I would separate the gathering and
reporting responsibilties into separate binaries to allow for more granular
control over deployments. This would prevent the two processes from
communicating via the DB (realistically they shouldn't be  anyways), this means
that either a JSON or (more likely) a gRPC interface would be established
allowing the two processes to run on different hosts. This would also allow
this service to scale horizontally for metrics gathering across a cluster or
multiple reporters behind a load balancer to handle higher load.

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
