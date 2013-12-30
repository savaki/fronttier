package mock

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseWriter struct {
	header     http.Header
	StatusCode int
	Content    []byte
	reader     io.ReadCloser
}

func (self *ResponseWriter) Header() http.Header {
	if self.header == nil {
		self.header = make(http.Header)
	}
	return self.header
}

func (self *ResponseWriter) WriteHeader(statusCode int) {
	self.StatusCode = statusCode
}

func (self *ResponseWriter) Write(data []byte) (int, error) {
	self.Content = append(self.Content, data...)
	return len(data), nil
}

func (self *ResponseWriter) WriteLater(reader io.ReadCloser) {
	self.reader = reader
}

func (self *ResponseWriter) WriteTo(writer io.Writer) {
	switch v := writer.(type) {
	case *ResponseWriter:
		v.reader = self.reader
		v.Content = self.Content

	default:
		if self.reader == nil {
			writer.Write(self.Content)
		} else {
			defer self.reader.Close()
			io.Copy(writer, self.reader)
		}
	}
}

func WriteTo(writer io.Writer, source io.ReadCloser) {
	switch v := writer.(type) {
	case *ResponseWriter:
		safeSource := source
		v.WriteLater(safeSource)

	default:
		io.Copy(writer, source)
	}
}

func (self *ResponseWriter) Bytes() []byte {
	buffer := &bytes.Buffer{}
	self.WriteTo(buffer)
	return buffer.Bytes()
}

func (self *ResponseWriter) String() string {
	return string(self.Bytes())
}
