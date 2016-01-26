package flags

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/alistanis/st/parse"
	"github.com/alistanis/st/sterrors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFlags(t *testing.T) {
	Convey("Flag testing", t, func() {

		Convey("We can set the mode to camel case", func() {
			ResetFlags()
			SetArgs([]string{"-c", ""})
			err := ParseFlags()
			So(err, ShouldBeNil)
			So(c, ShouldBeTrue)
		})

		Convey("We can set the mode to snake case", func() {
			ResetFlags()
			SetArgs([]string{"-s", ""})
			err := ParseFlags()
			So(err, ShouldBeNil)
			So(s, ShouldBeTrue)
		})

		Convey("Append mode is skip existing by default", func() {
			So(AppendMode, ShouldEqual, parse.SkipExisting)
		})

		Convey("We can set the append mode to overwrite", func() {
			ResetFlags()
			SetArgs([]string{"-o", ""})
			err := ParseFlags()
			So(err, ShouldBeNil)
			So(AppendMode, ShouldEqual, parse.Overwrite)
		})

		Convey("We can set append mode to append", func() {
			ResetFlags()
			SetArgs([]string{"-a", ""})
			err := ParseFlags()
			So(err, ShouldBeNil)
			So(AppendMode, ShouldEqual, parse.Append)
		})

		Convey("We can set ignored fields", func() {
			ResetFlags()
			SetArgs([]string{"-i", "ignore,this,field", ""})
			err := ParseFlags()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(parse.IgnoredFields, []string{"ignore", "this", "field"}), ShouldBeTrue)
		})

		Convey("We can set ignored structs", func() {
			ResetFlags()
			SetArgs([]string{"-is", "ignore,these,structs", ""})
			err := ParseFlags()
			So(err, ShouldBeNil)
			So(reflect.DeepEqual(parse.IgnoredStructs, []string{"ignore", "these", "structs"}), ShouldBeTrue)
		})

	})
}

func TestFlagErrors(t *testing.T) {
	Convey("Flag error testing", t, func() {

		Convey("We can return an error when no path is given", func() {
			ResetFlags()
			SetArgs([]string{})
			err := ParseFlags()
			So(err, ShouldNotBeNil)
			So(err, ShouldEqual, sterrors.NoPathsGiven)
		})

		Convey("Given a set of mismatched case flags", func() {
			Convey("A mutually exclusive parameters error is given", func() {
				ResetFlags()
				SetArgs([]string{"-c", "-s", ""})
				err := ParseFlags()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, sterrors.MutuallyExclusiveParameters("c", "s").Error())
			})
		})

		Convey("Given a set of mismatched append mode flags", func() {
			Convey("A mutually exclusive parameters error is given", func() {
				ResetFlags()
				SetArgs([]string{"-a", "-o", ""})
				err := ParseFlags()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, sterrors.MutuallyExclusiveParameters("o", "a").Error())
			})
		})

	})
}

func ResetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
}

func SetArgs(s []string) {
	os.Args = []string{os.Args[0]}
	os.Args = append(os.Args, s...)
}
