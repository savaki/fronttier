package sessions

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMemory(t *testing.T) {
	var store Store
	var session *Session
	var err error
	values := map[string]string{"hello": "world"}

	Convey("Given a Memory store", t, func() {
		store = Memory()

		Convey("When I create a session", func() {
			session, err = store.Create(values)

			Convey("Then I expect no errors", func() {
				So(err, ShouldBeNil)
				So(session, ShouldNotBeNil)
			})

			Convey("And I expect to be able to retrieve the session via #Get", func() {
				actual, err := store.Get(session.Id)

				So(err, ShouldBeNil)
				So(actual, ShouldNotBeNil)
				So(actual, ShouldResemble, session)
			})

			Convey("And I expect to be able to delete the session via #Delete", func() {
				err = store.Delete(session.Id)

				So(err, ShouldBeNil)

				actual, err := store.Get(session.Id)
				So(err, ShouldBeNil)
				So(actual, ShouldBeNil)
			})
		})
	})
}
