package sessions

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestRedis(t *testing.T) {
	if os.Getenv("REDIS_HOST") == "" {
		Convey("Redis tests are disabled", t, func() {
			Convey("To test with redis set REDIS_HOST to the redis addr as follows:", func() {
				Convey("export REDIS_HOST=localhost:6379", func() {
				})
			})
		})

	} else {
		Convey("#Set", t, func() {
			Convey("should set variables in redis", func() {
				store := Redis().Namespace("testing").Build()
				session := &Session{
					Values: map[string]string{"hello": "world"},
				}
				store.Set("hello", session)

				actual, err := store.Get("hello")

				So(err, ShouldBeNil)
				So(actual, ShouldNotBeNil)
				So(actual.Values["hello"], ShouldEqual, "world")
			})
		})

		Convey("#Create", t, func() {
			Convey("should create a new session with a new sessionId and persist it into redis", func() {
				store := Redis().Namespace("testing").Build()
				values := map[string]string{"foo": "bar"}

				// When
				session, err := store.Create(values)

				// Then
				So(err, ShouldBeNil)
				So(session, ShouldNotBeNil)

				actual, _ := store.Get(session.Id)
				So(actual, ShouldResemble, session)
				So(actual.Values["foo"], ShouldEqual, "bar")
			})
		})

		Convey("#Build", t, func() {
			Convey("should default the namespace to 'sessions'", func() {
				store := Redis().Build()
				r := store.(*redis)

				So(r.namespace, ShouldEqual, "sessions")
			})

			Convey("should assign #MaxIdle", func() {
				maxIdle := 7
				builder := Redis().MaxIdle(maxIdle)

				So(builder.maxIdle, ShouldEqual, maxIdle)
			})

			Convey("should assign #Host", func() {
				host := "localhost:1234"
				builder := Redis().Host(host)

				So(builder.host, ShouldEqual, host)
			})
		})

		Convey("#Delete", t, func() {
			Convey("should delete values from redis", func() {
				store := Redis().Namespace("testing").Build()
				values := map[string]string{"foo": "bar"}

				// When
				session, _ := store.Create(values)
				err := store.Delete(session.Id)

				// Then
				So(err, ShouldBeNil)
				So(session, ShouldNotBeNil)
				So(session.Id, ShouldNotEqual, "")

				value, err := store.Get(session.Id)
				So(err, ShouldBeNil)
				So(value, ShouldBeNil)
			})
		})
	}
}
