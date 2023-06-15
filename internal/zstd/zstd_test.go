/*
 *
 * Copyright 2021 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package zstd

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/test/bufconn"

	"github.com/mostynb/go-grpc-compression/internal/testserver"
)

const (
	bufSize = 1024
	message = "Message Request zSTD"
)

func TestRegisteredCompression(t *testing.T) {
	clobbering := true
	PretendInit(clobbering)

	for _, lvl := range []zstd.EncoderLevel{zstd.SpeedFastest, zstd.SpeedDefault, zstd.SpeedBetterCompression, zstd.SpeedBestCompression} {
		require.NoError(t, SetLevel(lvl))

		comp := encoding.GetCompressor(Name)
		require.NotNil(t, comp)
		assert.Equal(t, Name, comp.Name())

		buf := bytes.NewBuffer(make([]byte, 0, bufSize))
		wc, err := comp.Compress(buf)
		require.NoError(t, err)

		_, err = wc.Write([]byte(message))
		require.NoError(t, err)
		assert.NoError(t, wc.Close())

		r, err := comp.Decompress(buf)
		require.NoError(t, err)
		expected, err := ioutil.ReadAll(r)
		require.NoError(t, err)

		assert.Equal(t, message, string(expected))
	}
}

func TestRoundTrip(t *testing.T) {
	clobbering := true
	PretendInit(clobbering)

	lis := bufconn.Listen(bufSize)
	t.Cleanup(func() {
		assert.NoError(t, lis.Close())
	})

	done := make(chan struct{}, 1)

	s := grpc.NewServer()
	defer func() {
		s.GracefulStop()
		<-done
	}()
	testserver.RegisterTestServerServer(s, &testserver.EchoTestServer{})
	go func() {
		if err := s.Serve(lis); err != nil && err != grpc.ErrServerStopped {
			t.Errorf("Server exited with error: %v", err)
		}
		done <- struct{}{}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithDefaultCallOptions(grpc.UseCompressor(Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, conn.Close())
	})

	client := testserver.NewTestServerClient(conn)
	resp, err := client.SendMessage(ctx, &testserver.MessageRequest{Request: message})
	require.NoError(t, err)
	assert.Equal(t, message, resp.Response)
}
