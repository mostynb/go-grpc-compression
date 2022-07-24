# go-grpc-compression

This respository contains go gRPC encoding wrappers for some useful compression
algorithms that are not available in google.golang.org/grpc.

* https://github.com/golang/snappy
* https://github.com/klauspost/compress/tree/master/zstd
* https://github.com/pierrec/lz4

The following algorithms also have experimental implementations, which have
not been tested as much as those above. These may be changed significantly, or
even removed from this library at a future point.

* https://github.com/klauspost/compress/tree/master/s2
