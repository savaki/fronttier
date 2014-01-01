package sessions

import (
	"errors"
	redigo "github.com/garyburd/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type Command struct {
	Name string
	Args []interface{}
}

type MockConn struct {
	Done       []Command
	Sent       []Command
	Response   interface{}
	Error      error
	DoErr      error
	SendErr    error
	FlushErr   error
	ReceiveErr error
}

func (self *MockConn) ConnFactory() (redigo.Conn, error) {
	return self, nil
}

func (self *MockConn) Close() error {
	return self.Error
}

func (self *MockConn) Err() error {
	return self.Error
}

func (self *MockConn) Do(commandName string, args ...interface{}) (interface{}, error) {
	if self.DoErr != nil {
		return nil, self.DoErr
	}

	command := Command{commandName, []interface{}(args)}
	self.Done = append(self.Done, command)

	return self.Response, self.Error
}

func (self *MockConn) Send(commandName string, args ...interface{}) error {
	if self.SendErr != nil {
		return self.SendErr
	}

	command := Command{commandName, []interface{}(args)}
	self.Sent = append(self.Sent, command)

	return self.Error
}

func (self *MockConn) Flush() error {
	if self.FlushErr != nil {
		return self.FlushErr
	}
	return self.Error
}

func (self *MockConn) Receive() (interface{}, error) {
	if self.ReceiveErr != nil {
		return nil, self.ReceiveErr
	}
	return self.Response, self.Error
}

func TestMockConn(t *testing.T) {
	var mock *MockConn

	Convey("Given a MockConn", t, func() {
		mock = &MockConn{}

		Convey("MockConn should implement redis.Conn", func() {
			var conn redigo.Conn
			conn = mock

			So(conn, ShouldNotBeNil)
		})

		Convey("When I set MockConn#err", func() {
			mock.Error = errors.New("boom")

			Convey("Then I expect #Err to return the err", func() {
				So(mock.Err(), ShouldEqual, mock.Error)
			})
		})

		Convey("When I invoke #Do", func() {
			mock.Response = "argle"
			response, err := mock.Do("DEL", "abc")

			Convey("Then I expect the command to be captured", func() {
				So(len(mock.Done), ShouldEqual, 1)

				command := mock.Done[0]
				So(command.Name, ShouldEqual, "DEL")
				So(len(command.Args), ShouldEqual, 1)
				So(command.Args[0], ShouldEqual, "abc")
			})

			Convey("And I expect no errors", func() {
				So(err, ShouldBeNil)
			})

			Convey("And I expect the response to be returned", func() {
				So(response, ShouldEqual, mock.Response)
			})
		})

		Convey("When I invoke #Sent", func() {
			err := mock.Send("DEL", "abc")

			Convey("Then I expect the command to be captured", func() {
				So(len(mock.Sent), ShouldEqual, 1)

				command := mock.Sent[0]
				So(command.Name, ShouldEqual, "DEL")
				So(len(command.Args), ShouldEqual, 1)
				So(command.Args[0], ShouldEqual, "abc")
			})

			Convey("And I expect no errors", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
