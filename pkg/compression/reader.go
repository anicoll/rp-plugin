package compression

import (
	"github.com/klauspost/compress/zstd"
)

// Reader
type readCloser struct {
	zstdDecoder *zstd.Decoder
}

// NewReadCloser returns a stateless reader for decompression using zstd.
func NewReadCloser() (ReadCloser, error) {
	reader, err := zstd.NewReader(nil, zstd.WithDecoderConcurrency(0))
	if err != nil {
		return nil, err
	}
	return &readCloser{
		zstdDecoder: reader,
	}, nil
}

func (rc *readCloser) Close() error {
	rc.zstdDecoder.Close()
	return nil
}

func (rc *readCloser) DecodeAll(p []byte) ([]byte, error) {
	return rc.zstdDecoder.DecodeAll(p, nil)
}
