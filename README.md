# event-spread
Model news travelling in a map for e.g. more realistic game morality systems

TODO(cripplet): Add documentation in README.

## Testing

```
$ bazel --version
bazel 3.2.0

$ bazel test //lib/core:all --features=race --runs_per_test=100
```
## Running

```
$ bazel run //bin:main
```

## Client

See [grpc_cli](https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#grpc-cli).
Useful for testing.

*N.B.*: Ensure the `gflags` is installed: `sudo apt-get install -y libgflags-dev libgtest-dev libc++-dev clang`.

See https://github.com/grpc/grpc/issues/8582 ¯\\_(ツ)\_/¯

```
$ grpc_cli ls -l localhost:8080
$ grpc_cli call localhost:8080 EventSpreadService.AddEvent "event: {spread_rate: 0}"
connecting to localhost:8080
Rpc succeeded with OK status
```
