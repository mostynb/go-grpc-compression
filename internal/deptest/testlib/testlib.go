// Copyright 2020 Mostyn Bramley-Moore.
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

package testlib

import (
	"fmt"
	"io"

	"google.golang.org/grpc/encoding"
)

type Dummy string

var _ encoding.Compressor = Dummy("")

// AllNames lists every compressor in this repo, for testing.
var AllNames = []string{"zstd", "snappy", "lz4"}

func (d Dummy) Compress(io.Writer) (io.WriteCloser, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d Dummy) Decompress(r io.Reader) (io.Reader, error) {
	return nil, fmt.Errorf("not implemented")
}

func (d Dummy) Name() string {
	return string(d)
}

func init() {
	for _, name := range AllNames {
		// This test will not register the dummies in case
		// of existing registrations, this ensures the import
		// order of the deptest actually tests no-clobbering.
		if encoding.GetCompressor(name) == nil {
			encoding.RegisterCompressor(Dummy(name))
		}
	}
}
