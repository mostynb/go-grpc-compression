# go-grpc-compression

This respository contains go gRPC encoding wrappers for some useful compression
algorithms that are not available in google.golang.org/grpc.

* github.com/mostynb/go-grpc-compression/lz4 - using https://github.com/pierrec/lz4
* github.com/mostynb/go-grpc-compression/snappy - using https://github.com/golang/snappy
* github.com/mostynb/go-grpc-compression/zstd - using https://github.com/klauspost/compress/tree/master/zstd

The following algorithms also have experimental implementations, which have
not been tested as much as those above. These may be changed significantly, or
even removed from this library at a future point.

* github.com/mostynb/go-grpc-compression/experimental/klauspost_snappy - using https://github.com/klauspost/compress/tree/master/s2
  in snappy compatibility mode
* github.com/mostynb/go-grpc-compression/experimental/s2 - using https://github.com/klauspost/compress/tree/master/s2

Importing any of the packages above will override any previously registered
encoders with the same name. If you would prefer imports to only register
the encoder if there is no previously registered encoder with the same name,
then you should instead import one of the following packages:

* github.com/mostynb/go-grpc-compression/nonclobbering/lz4
* github.com/mostynb/go-grpc-compression/nonclobbering/snappy
* github.com/mostynb/go-grpc-compression/nonclobbering/zstd
* github.com/mostynb/go-grpc-compression/nonclobbering/experimental/klauspost_snappy
* github.com/mostynb/go-grpc-compression/nonclobbering/experimental/s2
