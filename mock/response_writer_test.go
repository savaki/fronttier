package mock

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"testing"
)

func TestResponseWriter(t *testing.T) {
	Convey("#Header", t, func() {
		Convey("should capture headers", func() {
			w := &ResponseWriter{}

			// When
			w.Header().Set("Hello", "world")

			// Then
			So(w.Header().Get("Hello"), ShouldEqual, "world")
		})
	})

	Convey("#WriteHeader", t, func() {
		Convey("should capture status code", func() {
			w := &ResponseWriter{}
			statusCode := 201

			// When
			w.WriteHeader(statusCode)

			// Then
			So(w.StatusCode, ShouldEqual, statusCode)
		})
	})

	Convey("#Write", t, func() {
		Convey("should save the content", func() {
			w := &ResponseWriter{}
			content := []byte("hello world")

			// When
			w.Write(content)

			// Then
			So(string(w.Content), ShouldEqual, string(content))
		})
	})

	Convey("#WriteLater", t, func() {
		var w *ResponseWriter

		Convey("Read should only be invoked when the thingie is called", func() {
			contents := "hello world"
			r := &MockReaderCloser{
				reader: bytes.NewReader([]byte(contents)),
			}
			w = &ResponseWriter{}
			w.WriteLater(r)

			So(r.readCount, ShouldEqual, 0)
			So(r.closeCount, ShouldEqual, 0)
			So(w.String(), ShouldEqual, contents)
			So(r.readCount, ShouldBeGreaterThan, 0)
			So(r.closeCount, ShouldBeGreaterThan, 0)
		})
	})

	Convey("#WriteTo", t, func() {
		contents := "hello world"

		Convey("Calls #WriteLater if target is a ResponseWriter", func() {
			r := &MockReaderCloser{data: []byte(contents)}
			w := &ResponseWriter{}

			WriteTo(w, r)

			So(len(w.Content), ShouldEqual, 0)
			So(w.reader, ShouldNotBeNil)
			So(w.String(), ShouldEqual, contents)
		})

		Convey("Copies data from source to target", func() {
			r := &MockReaderCloser{data: []byte(contents)}
			buffer := &bytes.Buffer{}

			WriteTo(buffer, r)

			So(buffer.String(), ShouldEqual, contents)
		})
	})

	Convey("#Bytes", t, func() {
		var w *ResponseWriter
		contents := "hello world"

		Convey("When []byte are written via #Write", func() {
			w = &ResponseWriter{}
			w.Write([]byte(contents))

			Convey("Then #Bytes should return those bytes", func() {
				So(string(w.Bytes()), ShouldEqual, contents)
			})

			Convey("Then #String should return the string", func() {
				So(w.String(), ShouldEqual, contents)
			})

			Convey("Then we CAN NOT call String() multiple times", func() {
				So(w.String(), ShouldEqual, contents)
			})
		})

		Convey("When []byte are written via #WriteLater", func() {
			var r io.ReadCloser
			data := []byte("hello world")

			Convey("Then #WriteTo should write to a io.Buffer", func() {
				r = &MockReaderCloser{data: data}
				a := &ResponseWriter{}
				a.WriteLater(r)

				buffer := &bytes.Buffer{}
				a.WriteTo(buffer)

				So(buffer.String(), ShouldEqual, string(data))
			})

			Convey("Then #WriteTo should write to another ResponseWriter", func() {
				r = &MockReaderCloser{data: data}
				a := &ResponseWriter{}
				a.WriteLater(r)

				b := &ResponseWriter{}
				a.WriteTo(b)

				So(len(b.Content), ShouldEqual, 0)
				So(b.reader, ShouldEqual, a.reader)

				buffer := &bytes.Buffer{}
				b.WriteTo(buffer)
				So(buffer.String(), ShouldEqual, string(data))
			})
		})
	})
}

type MockReaderCloser struct {
	readCount  int
	closeCount int
	data       []byte
	reader     io.Reader
}

func (self *MockReaderCloser) Read(p []byte) (int, error) {
	if self.reader == nil {
		self.reader = bytes.NewReader(self.data)
	}

	self.readCount++
	return self.reader.Read(p)
}

func (self *MockReaderCloser) Close() error {
	self.closeCount++
	return nil
}
