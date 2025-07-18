package main

import (
	"fmt"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAppCredsProcessing(t *testing.T) {
	testCreds := "./fixtures/credentials.yml.test"
	testData := []byte("data: 12345678987654321\n")

	err := os.WriteFile(testCreds, []byte{}, 0644)
	if err != nil {
		fmt.Printf("Error writing credentials file: %v\n", err)
		t.Fail()
	}
	defer func() {
		rmErr := os.Remove(testCreds)
		if rmErr != nil {
			fmt.Printf("Error removing test file: %s\n", rmErr.Error())
		}
	}()

	app := &App{
		opts: &AppOptions{
			CredsPath: testCreds,
			KeyFile:   "./fixtures/master.key",
		},
	}
	app.decodeEncKey()

	Convey("Writes correctly decoded credentials", t, func() {
		err = app.encodeFile(testData)

		Convey("Then it should not return an error", func() {
			So(err, ShouldBeNil)
		})
	})

	Convey("Read encoded credentials with correct key", t, func() {
		data, err := app.decodeFile()

		Convey("Then it should not return an error", func() {
			So(err, ShouldBeNil)
		})

		Convey("Then it decodes file correctly", func() {
			So(data, ShouldEqual, testData)
		})
	})

	Convey("Read encoded credentials with incorrect key", t, func() {
		app.Key = []byte{0x00, 0x01}
		_, err = app.decodeFile()

		Convey("Then it should return an error", func() {
			So(err, ShouldNotBeNil)
		})
	})
}
