package sessions

import (
	"encoding/json"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRedisWithMocks(t *testing.T) {
	var mock *MockConn
	var store Store
	var err error
	var s *Session
	var data []byte
	session := &Session{Id: "1234", Values: map[string]string{"hello": "world"}}
	key := "das key"

	Convey("Given a redis store with a mock conn", t, func() {
		data, _ = json.Marshal(session)
		mock = &MockConn{}
		store = Redis().ConnFactory(mock.ConnFactory).Build()

		Convey("When I call #Create", func() {
			s, err = store.Create(session.Values)

			Convey("Then I expect no error to be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then I expect a new session to be created", func() {
				So(s.Id, ShouldNotEqual, session.Id)
				So(s.Values, ShouldResemble, session.Values)
			})

			Convey("Then I expect Sent to have captured the command", func() {
				So(len(mock.Sent), ShouldEqual, 1)

				command := mock.Sent[0]
				So(command.Name, ShouldEqual, "SET")
				So(len(command.Args), ShouldEqual, 2)
				So(command.Args[0], ShouldContainSubstring, s.Id)

				data, _ = json.Marshal(s)
				So(command.Args[1], ShouldEqual, string(data))
			})
		})

		Convey("When I call #Set", func() {
			err = store.Set(key, session)

			Convey("Then I expect no error to be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then I expect Sent to have captured the command", func() {
				So(len(mock.Sent), ShouldEqual, 1)

				command := mock.Sent[0]
				So(command.Name, ShouldEqual, "SET")
				So(len(command.Args), ShouldEqual, 2)
				So(command.Args[0], ShouldContainSubstring, key)
				So(command.Args[1], ShouldEqual, string(data))
			})

			Convey("When I receive a #SendErr", func() {
				mock.SendErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					err = store.Set(key, session)

					So(err, ShouldEqual, mock.SendErr)
				})
			})

			Convey("When I receive a #FlushErr", func() {
				mock.FlushErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					err = store.Set(key, session)

					So(err, ShouldEqual, mock.FlushErr)
				})
			})

			Convey("When I receive a #ReceiveErr", func() {
				mock.ReceiveErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					err = store.Set(key, session)

					So(err, ShouldEqual, mock.ReceiveErr)
				})
			})
		})

		Convey("When I call #Get", func() {
			mock.Response = data
			value, err := store.Get(key)

			Convey("Then I expect no error to be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then I expect #Get to return the response", func() {
				So(value, ShouldResemble, session)
			})

			Convey("Then I expect Sent to have captured the command", func() {
				So(len(mock.Sent), ShouldEqual, 1)

				command := mock.Sent[0]
				So(command.Name, ShouldEqual, "GET")
				So(len(command.Args), ShouldEqual, 1)
				So(command.Args[0], ShouldContainSubstring, key)
			})

			Convey("When I receive a #SendErr", func() {
				mock.SendErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					_, err = store.Get(key)

					So(err, ShouldEqual, mock.SendErr)
				})
			})

			Convey("When I receive a #FlushErr", func() {
				mock.FlushErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					_, err = store.Get(key)

					So(err, ShouldEqual, mock.FlushErr)
				})
			})

			Convey("When I receive a #ReceiveErr", func() {
				mock.ReceiveErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					_, err = store.Get(key)

					So(err, ShouldEqual, mock.ReceiveErr)
				})
			})

			Convey("When #Receive returns nil", func() {
				mock.Response = nil

				Convey("Then I expect #Set to return an err", func() {
					v, err := store.Get(key)

					So(err, ShouldBeNil)
					So(v, ShouldBeNil)
				})
			})
		})

		Convey("When I call #Delete", func() {
			mock.Response = data
			err := store.Delete(key)

			Convey("Then I expect no error to be returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then I expect Sent to have captured the command", func() {
				So(len(mock.Sent), ShouldEqual, 1)

				command := mock.Sent[0]
				So(command.Name, ShouldEqual, "DEL")
				So(len(command.Args), ShouldEqual, 1)
				So(command.Args[0], ShouldContainSubstring, key)
			})

			Convey("When I receive a #SendErr", func() {
				mock.SendErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					err = store.Delete(key)

					So(err, ShouldEqual, mock.SendErr)
				})
			})

			Convey("When I receive a #FlushErr", func() {
				mock.FlushErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					err = store.Delete(key)

					So(err, ShouldEqual, mock.FlushErr)
				})
			})

			Convey("When I receive a #ReceiveErr", func() {
				mock.ReceiveErr = errors.New("boom")

				Convey("Then I expect #Set to return an err", func() {
					err = store.Delete(key)

					So(err, ShouldEqual, mock.ReceiveErr)
				})
			})
		})
	})
}
