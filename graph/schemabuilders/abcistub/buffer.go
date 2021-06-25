package abcistub

import (
	"bytes"
	"net/http"
)

// double as http.ResponseWriter and runtime.ClientResponse
type ResponseBuffer struct {
	code   int
	header http.Header
	buf    *bytes.Buffer
}

func (w *ResponseBuffer) Write(b []byte) (int, error) {
	w.buf = bytes.NewBuffer(b)
	return len(b), nil
}

func (w *ResponseBuffer) WriteHeader(statusCode int) {
	w.code = statusCode
}
func (w *ResponseBuffer) Header() http.Header {
	w.header = http.Header{}
	return w.header
}
func (w *ResponseBuffer) Finalize() []byte {
	return w.buf.Bytes()
}
func (w *ResponseBuffer) FinalizeString() string {
	return string(w.buf.Bytes())
}

func (w *ResponseBuffer) Code() int {
	return w.code
}

func (w *ResponseBuffer) Message() string {
	return ""
}

func (w *ResponseBuffer) GetHeader(key string) string {
	return ""
}

func (w *ResponseBuffer) Body() ReaderCloser {
	return ReaderCloser{
		buf: w.buf,
	}
}

type ReaderCloser struct {
	buf *bytes.Buffer
}

func (rc ReaderCloser) Clone() ReaderCloser {
	return ReaderCloser{
		buf: bytes.NewBuffer(rc.buf.Bytes()),
	}
}

func (rc ReaderCloser) Read(p []byte) (n int, err error) {
	return rc.buf.Read(p)
}

func (rc ReaderCloser) Close() error {
	return nil
}
