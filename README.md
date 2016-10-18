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
