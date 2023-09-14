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

// Package github.com/mostynb/go-grpc-compression/lz4 is a wrapper for
// using github.com/pierrec/lz4 with gRPC.
//
// If you import this package, it will register itself as the encoder for
// the "lz4" compressor, overriding any previously registered compressors
// with this name.
//
// If you don't want to override previously registered "lz4" compressors,
// then you should instead import
// github.com/mostynb/go-grpc-compression/nonclobbering/lz4
package lz4

import (
	internallz4 "github.com/mostynb/go-grpc-compression/internal/lz4"
)

const Name = internallz4.Name

func init() {
	clobbering := true
	internallz4.PretendInit(clobbering)
}
