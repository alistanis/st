package sterrors

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	con "github.com/smartystreets/goconvey/convey"
)

var (
	tempDir string
)

func init() {
	var err error
	tempDir, err = ioutil.TempDir("", "sterrors")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func TestErrors(t *testing.T) {
	con.Convey("Verbose printing tests", t, func() {
		Verbose = true
		stdout := os.Stdout
		fname := filepath.Join(tempDir, "stdout1")
		temp, err := os.Create(fname)
		con.So(err, con.ShouldBeNil)

		os.Stdout = temp
		testString := "This is a test string"
		Printf(testString)
		err = temp.Close()
		con.So(err, con.ShouldBeNil)

		output, err := ioutil.ReadFile(fname)
		con.Convey("We can print to stdout when verbose is enabled", func() {
			con.So(err, con.ShouldBeNil)
			con.So(string(output), con.ShouldEqual, testString)
		})

		Verbose = false

		fname = filepath.Join(tempDir, "stdout1")
		temp, err = os.Create(fname)
		con.So(err, con.ShouldBeNil)
		os.Stdout = temp
		Printf(testString)
		err = temp.Close()
		con.So(err, con.ShouldBeNil)
		output, err = ioutil.ReadFile(fname)

		con.Convey("We can not print to stdout when verbose is enabled", func() {
			con.So(err, con.ShouldBeNil)
			con.So(string(output), con.ShouldEqual, "")
		})
		os.Stdout = stdout
	})

	con.Convey("Mutually exclusive parameters returns an error in the format we expect", t, func() {
		err := ErrMutuallyExclusiveParameters("1", "2")
		con.So(err.Error(), con.ShouldEqual, "Mutually exclusive parameters provided: 1 and 2")
	})

	con.Convey("We can test http formatting", t, func() {
		testErr := errors.New("Test error")
		errBytes := FormatHTTPError(testErr, 400)
		type ErrResp struct {
			Message string `json:"error"`
			Code    int    `json:"status_code"`
		}
		fmt.Println(string(errBytes))
		errResp := &ErrResp{}
		err := json.Unmarshal(errBytes, &errResp)
		con.So(err, con.ShouldBeNil)
		con.So(errResp.Code, con.ShouldEqual, 400)
		con.So(errResp.Message, con.ShouldEqual, testErr.Error())
	})

}
