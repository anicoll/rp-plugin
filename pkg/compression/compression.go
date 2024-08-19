package compression

type WriteCloser interface {
	Write(p []byte) (n int, err error)
	Close() error
}

type ReadCloser interface {
	Close() error
	DecodeAll(p []byte) ([]byte, error)
}
