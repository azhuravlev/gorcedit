package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMarshalRubyString(t *testing.T) {
	Convey("Given a plain string", t, func() {
		input := "hello world"
		data := []byte(input)

		Convey("When marshalled and then unmarshalled", func() {
			marshalled, err := MarshalRubyString(data)
			So(err, ShouldBeNil)

			unmarshalled, err := UnmarshalRubyString(marshalled)
			So(err, ShouldBeNil)

			Convey("Then the original string should be preserved", func() {
				So(string(unmarshalled), ShouldEqual, input)
			})
		})
	})

	Convey("Given data without Ruby Marshal header", t, func() {
		raw := []byte("raw plaintext")

		Convey("When unmarshalled", func() {
			_, err := UnmarshalRubyString(raw)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a byte array with unsupported tag", t, func() {
		unsupported := []byte{0x04, 0x08, 0x30}

		Convey("When unmarshalled", func() {
			_, err := UnmarshalRubyString(unsupported)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a byte array with marshalled text", t, func() {
		unsupported := []byte{0x04, 0x08, 0x22, 0x06, 0xAA}

		Convey("When unmarshalled", func() {
			result, err := UnmarshalRubyString(unsupported)

			Convey("Then it should not return an error", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then it returns decoded string", func() {
				So(result, ShouldEqual, []byte{0xAA})
			})
		})
	})
}
