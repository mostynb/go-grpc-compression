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

// package deptest ensures that the go-grpc-compression sub-packages will
// not clobber an existing registration.
package deptest

import (
	"testing"

	// This import happens first for the test, keep it ahead of
	// the three unnamed imports below.
	"github.com/mostynb/go-grpc-compression/internal/deptest/testlib"

	// If these were moved above, the test would fail because
	// testlib has the same no-clobber logic as the main packages.
	_ "github.com/mostynb/go-grpc-compression/lz4"
	_ "github.com/mostynb/go-grpc-compression/snappy"
	_ "github.com/mostynb/go-grpc-compression/zstd"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/encoding"
)

func TestCompressorNotClobbered(t *testing.T) {
	// Because the deptest/lib imports first, it's init() function
	// registers dummy compressors. The following libraries do not
	// clobber, so we should find Dummy compressors.
	for _, name := range testlib.AllNames {
		_, ok := encoding.GetCompressor(name).(testlib.Dummy)
		require.True(t, ok)
	}
}
