package parse

import (
	"os"
	"reflect"
	"testing"

	"github.com/alistanis/st/sterrors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFlags(t *testing.T) {
	Convey("Flag testing", t, func() {
		// We trick the flag parser into thinking there is an additional parameter, when really, there isn't. This allows
		// us to bypass the check in ParseFlags() for a path, but this will be caught in parse or in main either way
		Convey("We can set the mode to camel case", func() {
			SetArgs([]string{"-c", ""})
			err := Flags()
			So(err, ShouldBeNil)
			So(c, ShouldBeTrue)
		})

		Convey("We can set the mode to snake case", func() {
			SetArgs([]string{"-s", ""})
			err := Flags()
			So(err, ShouldBeNil)
			So(s, ShouldBeTrue)
		})

		Convey("Append mode is skip existing by default", func() {
			So(AppendMode, ShouldEqual, SkipExisting)
		})

		Convey("We can set the append mode to overwrite", func() {
			SetArgs([]string{"-o", ""})
			err := Flags()
			So(err, ShouldBeNil)
			So(AppendMode, ShouldEqual, Overwrite)
		})

		Convey("We can set append mode to append", func() {
			SetArgs([]string{"-a", ""})
			err := Flags()
			So(err, ShouldBeNil)
			So(AppendMode, ShouldEqual, Append)
		})

		Convey("We can set ignored fields", func() {
			SetArgs([]string{"-i", "ignore,this,field", ""})
			err := Flags()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(IgnoredFields, []string{"ignore", "this", "field"}), ShouldBeTrue)
		})

		Convey("We can set ignored structs", func() {
			SetArgs([]string{"-is", "ignore,these,structs", ""})
			err := Flags()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(IgnoredStructs, []string{"ignore", "these", "structs"}), ShouldBeTrue)
		})

	})
}

func TestGoGenerateDetection(t *testing.T) {
	Convey("Go generate detection testing", t, func() {
		Convey("Given a fake $GOFILE", func() {

			Convey("We can trick make sure that flags are set properly", func() {
				os.Setenv("GOFILE", "test.go")
				SetArgs([]string{})
				err := Flags()
				So(err, ShouldBeNil)
				So(GoFile, ShouldEqual, "test.go")
			})
			Reset(func() {
				os.Setenv("GOFILE", "")
			})
		})
	})
}

func TestFlagErrors(t *testing.T) {
	Convey("Flag error testing", t, func() {

		Convey("We can return an error when no path is given", func() {
			SetArgs([]string{})
			err := Flags()
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, sterrors.ErrNoPathsGiven)
		})

		Convey("Given a set of mismatched case flags", func() {
			Convey("A mutually exclusive parameters error is given", func() {
				SetArgs([]string{"-c", "-s", ""})
				err := Flags()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, sterrors.ErrMutuallyExclusiveParameters("c", "s").Error())
			})
		})

		Convey("Given a set of mismatched append mode flags", func() {
			Convey("A mutually exclusive parameters error is given", func() {
				SetArgs([]string{"-a", "-o", ""})
				err := Flags()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, sterrors.ErrMutuallyExclusiveParameters("o", "a").Error())
			})
		})

	})
}
