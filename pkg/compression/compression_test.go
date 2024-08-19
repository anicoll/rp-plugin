package compression

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func FuzzCompressDecompress(f *testing.F) {
	rc, err := NewReadCloser()
	assert.NoError(f, err)
	testcases := []string{"Hello, world", " ", "!12345"}
	for _, tc := range testcases {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}
	f.Fuzz(func(t *testing.T, orig string) {
		out := new(bytes.Buffer)
		wc, err := NewWriteCloser(out)
		assert.NoError(t, err)
		_, err = wc.Write([]byte(orig))
		assert.NoError(t, err)
		assert.NoError(t, wc.Close())

		decoded, err := rc.DecodeAll(out.Bytes())
		assert.NoError(t, err)
		decodedString := string(decoded)
		if utf8.ValidString(orig) && !utf8.ValidString(decodedString) {
			t.Errorf("Reverse produced invalid UTF-8 string %q", decodedString)
		}
	})
}
