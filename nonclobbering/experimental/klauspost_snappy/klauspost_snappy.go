// Copyright 2023 Mostyn Bramley-Moore.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package github.com/mostynb/go-grpc-compression/nonclobbering/klauspost_snappy
// is a wrapper for using github.com/klauspost/compress/s2 in snappy
// compatibility mode with gRPC. It might be more efficient than
// github.com/mostynb/go-grpc-compression/snappy and
// github.com/mostynb/go-grpc-compression/nonclobbering/snappy
// packages which make use of github.com/golang/snappy.
//
// Note that this is registered under the name "snappy" with gRPC, so only
// one of these packages should be used at a time.
//
// If you import this package, it will only register itself as the encoder
// for the "snappy" compressor if no other compressors have already been
// registered with that name.
//
// If you do want to override previously registered "snappy" compressors,
// then you should instead import
// github.com/mostynb/go-grpc-compression/klauspost_snappy
package klauspost_snappy

import (
	internalsnappy "github.com/mostynb/go-grpc-compression/internal/klauspost_snappy"
)

const Name = internalsnappy.Name

func init() {
	clobbering := false
	internalsnappy.PretendInit(clobbering)
}
