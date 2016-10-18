# Metrics gathering and reporting

This has only been tested on Linux, and depends on Docker to build.

## Building

Run `make` or `make build` to compile your app.  This will use a Docker image
to build your app, with the current directory volume-mounted into place.  This
will store incremental state for the fastest possible build.  Run `make
all-build` to build for all architectures.

Run `make clean` to clean up.

### Notes

* Only the amd64 arch for darwin and linux have been fully tested.
* Other OSs and ARCHs build successfully but have not been tested.
* darwin/arm fails to build. I do not fully understand why at this time.
