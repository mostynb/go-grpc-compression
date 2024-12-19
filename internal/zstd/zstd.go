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

// Package zstd is a wrapper for using github.com/klauspost/compress/zstd
// with gRPC.
package zstd

import (
	"errors"
	"io"
	"runtime"
	"sync"

	"github.com/klauspost/compress/zstd"
	"google.golang.org/grpc/encoding"
)

const Name = "zstd"

var encoderOptions = []zstd.EOption{
	// The default zstd window size is 8MB, which is much larger than the
	// typical RPC message and wastes a bunch of memory.
	zstd.WithWindowSize(512 * 1024),
}

var decoderOptions = []zstd.DOption{
	// If the decoder concurrency level is not 1, we would need to call
	// Close() to avoid leaking resources when the object is released
	// from compressor.decoderPool.
	zstd.WithDecoderConcurrency(1),
}

// We will set a finalizer on these objects, so when the go-grpc code is
// finished with them, they will be added back to compressor.decoderPool.
type decoderWrapper struct {
	*zstd.Decoder
}

type encoderWrapper struct {
	*zstd.Encoder
	pool *sync.Pool
}

type compressor struct {
	encoderPool sync.Pool
	decoderPool sync.Pool
}

func PretendInit(clobbering bool) {
	if !clobbering && encoding.GetCompressor(Name) != nil {
		return
	}

	encoding.RegisterCompressor(&compressor{})
}

var ErrNotInUse = errors.New("SetLevel ineffective because another zstd compressor has been registered")

// SetLevel updates the registered compressor to use a particular compression
// level. NOTE: this function must only be called from an init function, and
// is not threadsafe.
func SetLevel(level zstd.EncoderLevel) error {
	_, ok := encoding.GetCompressor(Name).(*compressor)
	if !ok {
		return ErrNotInUse
	}

	encoderOptions = append(encoderOptions, zstd.WithEncoderLevel(level))
	return nil
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	var err error
	var found bool
	var encoder *zstd.Encoder

	encoder, found = c.encoderPool.Get().(*zstd.Encoder)
	if !found {
		encoder, err = zstd.NewWriter(w, encoderOptions...)
		if err != nil {
			return nil, err
		}
	} else {
		encoder.Reset(w)
	}

	wrapper := &encoderWrapper{Encoder: encoder, pool: &c.encoderPool}
	runtime.SetFinalizer(wrapper, func(ew *encoderWrapper) {
		ew.Reset(nil)
		c.encoderPool.Put(ew.Encoder)
	})

	return wrapper, nil
}

func (w *encoderWrapper) Close() error {
	err := w.Encoder.Close()
	w.pool.Put(w.Encoder)
	return err
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	var err error
	var found bool
	var decoder *zstd.Decoder

	// Note: avoid the use of zstd.Decoder.DecodeAll here, since
	// malicious payloads could DoS us with a decompression bomb.

	decoder, found = c.decoderPool.Get().(*zstd.Decoder)
	if !found {
		decoder, err = zstd.NewReader(r, decoderOptions...)
		if err != nil {
			return nil, err
		}
	} else {
		err = decoder.Reset(r)
		if err != nil {
			c.decoderPool.Put(decoder)
			return nil, err
		}
	}

	wrapper := &decoderWrapper{Decoder: decoder}
	runtime.SetFinalizer(wrapper, func(dw *decoderWrapper) {
		err := dw.Reset(nil)
		if err == nil {
			c.decoderPool.Put(dw.Decoder)
		}
	})

	return wrapper, nil
}

func (c *compressor) Name() string {
	return Name
}
