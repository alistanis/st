package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/alistanis/st/parse"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	testData = `package test
type TestStruct struct {
	Field string
}
`
	expectedWrittenData = strings.Replace(`package test

type TestStruct struct {
	Field string %sjson:"field"%s
}
`, "%s", "`", -1)
	tempDir  string
	lastExit int
)

func init() {
	var err error
	tempDir, err = ioutil.TempDir("", "st_main_test_")
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func TestRun(t *testing.T) {
	Convey("Given a main program", t, func() {

		exitFunction = func(code int) {
			lastExit = code
		}
		Convey("We can hijack the exit function and test it", func() {
			exit(5)
			So(lastExit, ShouldEqual, 5)
		})

		Convey("We can call the run function", func() {
			// bury stderr so we don't see it; stdout is small here so we'll leave it
			tempStderr, err := ioutil.TempFile(tempDir, "stderr")
			So(err, ShouldBeNil)
			oldStderr := os.Stderr
			os.Stderr = tempStderr
			i := run()
			os.Stderr = oldStderr
			err = tempStderr.Close()
			Convey("without command line args it returns -1", func() {
				So(i, ShouldEqual, -1)
				So(lastExit, ShouldEqual, -2)
			})
		})
		Convey("Given a temporary file", func() {
			f, err := ioutil.TempFile(tempDir, "")
			So(err, ShouldBeNil)

			parse.SetArgs([]string{"-s", ""})
			i := run()

			So(err, ShouldBeNil)
			So(i, ShouldEqual, -1)

			f.WriteString(testData)
			parse.SetArgs([]string{"-s", "-w", f.Name()})
			i = run()
			So(i, ShouldEqual, 0)
			data, err := ioutil.ReadFile(f.Name())
			So(err, ShouldBeNil)
			So(string(data), ShouldEqual, expectedWrittenData)
		})
	})
}

// This doesn't show up in GoConvey's test coverage, but it actually does allow for full and complete coverage
// os.Exit() is not a hookable call, and goconvey isn't sophisticated enough to detect this kind of pattern yet
// Essentially what we do here is write our data to a file, pass that data as arguments to a new process with the
// additional -test.run=TestCrashingMain flag, and set an environment variable for that process, CRASH=1, so that
// when it runs, it will set the exitFunction (so that the default is used) and then run main. We then check the return
// and exit status of the main call, and we make sure that it was successful. This way, if we alter main() and we screw
// it up, we'll catch this when we run the tests
func TestCrashingMain(t *testing.T) {
	Convey("We can fake running main", t, func() {

		if os.Getenv("CRASH") == "1" {
			exitFunction = nil
			main()
		}

		f, err := ioutil.TempFile(tempDir, "")
		So(err, ShouldBeNil)
		f.WriteString(testData)
		args := "-test.run=TestCrashingMain -s " + f.Name()
		cmd := exec.Command(os.Args[0], args)
		cmd.Env = append(os.Environ(), "CRASH=1")

		_, err = cmd.Output()
		fail := false
		if e, ok := err.(*exec.ExitError); ok && !e.Success() {
			fail = true
		}
		So(fail, ShouldBeFalse)
	})
}
