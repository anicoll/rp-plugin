package compression

import (
	"bytes"

	"github.com/klauspost/compress/zstd"
)

// Writer

type writeCloser struct {
	zstdEncoder *zstd.Encoder
}

func NewWriteCloser(w *bytes.Buffer) (WriteCloser, error) {
	writer, err := zstd.NewWriter(w)
	if err != nil {
		return nil, err
	}
	return &writeCloser{
		zstdEncoder: writer,
	}, nil
}

func (wc *writeCloser) Close() error {
	return wc.zstdEncoder.Close()
}

func (wc *writeCloser) Write(p []byte) (n int, err error) {
	return wc.zstdEncoder.Write(p)
}
