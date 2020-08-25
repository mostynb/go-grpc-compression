/*
 *
 * Copyright 2017 gRPC authors.
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

// Package zstd is a wrapper for using github.com/klauspost/compress/zstd
// with gRPC.
package zstd

// This code is based upon the gzip wrapper in github.com/grpc/grpc-go:
// https://github.com/grpc/grpc-go/blob/master/encoding/gzip/gzip.go

import (
	"io"
	"io/ioutil"
	"runtime"
	"sync"

	zstdlib "github.com/klauspost/compress/zstd"
	"google.golang.org/grpc/encoding"
)

const Name = "zstd"

type compressor struct {
	poolCompressor   sync.Pool
	poolDecompressor sync.Pool
}

type writer struct {
	*zstdlib.Encoder
	pool *sync.Pool
}

type reader struct {
	*zstdlib.Decoder
	pool *sync.Pool
}

func init() {
	c := &compressor{}
	c.poolCompressor.New = func() interface{} {
		w, err := zstdlib.NewWriter(ioutil.Discard)
		if err != nil {
			panic(err)
		}
		writer := &writer{Encoder: w, pool: &c.poolCompressor}
		runtime.SetFinalizer(writer, finalizeWriter)
		return writer
	}
	encoding.RegisterCompressor(c)
}

// SetLevel updates the registered compressor to use a particular compression
// level. NOTE: this function must only be called from an init function, and
// is not threadsafe.
func SetLevel(level zstdlib.EncoderLevel) error {
	c := encoding.GetCompressor(Name).(*compressor)
	c.poolCompressor.New = func() interface{} {
		w, err := zstdlib.NewWriter(ioutil.Discard,
			zstdlib.WithEncoderLevel(level))
		if err != nil {
			return err
		}

		writer := &writer{Encoder: w, pool: &c.poolCompressor}
		runtime.SetFinalizer(writer, finalizeWriter)
		return writer
	}

	return nil
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	z := c.poolCompressor.Get().(*writer)
	z.Encoder.Reset(w)
	return z, nil
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	z, inPool := c.poolDecompressor.Get().(*reader)
	if !inPool {
		newZ, err := zstdlib.NewReader(r)
		if err != nil {
			return nil, err
		}
		reader := &reader{Decoder: newZ, pool: &c.poolDecompressor}
		runtime.SetFinalizer(reader, finalizeReader)
		return reader, nil
	}
	if err := z.Reset(r); err != nil {
		c.poolDecompressor.Put(z)
		return nil, err
	}
	return z, nil
}

func (c *compressor) Name() string {
	return Name
}

func (z *writer) Close() error {
	err := z.Encoder.Close()
	z.pool.Put(z)
	return err
}

func (z *reader) Read(p []byte) (n int, err error) {
	n, err = z.Decoder.Read(p)
	if err == io.EOF {
		z.pool.Put(z)
	}
	return n, err
}

func finalizeReader(r *reader) {
	if r.Decoder != nil {
		r.Decoder.Close()
	}
}

func finalizeWriter(w *writer) {
	if w.Encoder != nil {
		w.Encoder.Close()
	}
}
