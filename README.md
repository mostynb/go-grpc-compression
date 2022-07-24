# go-grpc-compression

This respository contains go gRPC encoding wrappers for some useful compression
algorithms that are not available in google.golang.org/grpc.

* snappy - using https://github.com/golang/snappy
* zstd - using https://github.com/klauspost/compress/tree/master/zstd
* lz4 - using https://github.com/pierrec/lz4

The following algorithms also have experimental implementations, which have
not been tested as much as those above. These may be changed significantly, or
even removed from this library at a future point.

* experimental/s2 - using https://github.com/klauspost/compress/tree/master/s2
* experimental/klauspost_snappy - using https://github.com/klauspost/compress/tree/master/s2
  in snappy compatibility mode
