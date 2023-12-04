package errhandle

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)


var timeReg = regexp.MustCompile(`^[\d]{4}-[\d]{2}-[\d]{2}T[\d]{2}:[\d]{2}:[\d]{2}\.[\d]{3}Z$`)

func TestLogGenerate(t *testing.T) {
	mnts := []string{"", "./testdata/"}
	files := []string{"./testdata/allfile.log", "mntfile.log"}

	for i, f := range files {
		l := LogInit("file:"+f, mnts[i])
		l.PrintErrorMessage("", fmt.Errorf(""), false)
		l.LogFlush()

		target := filepath.Join(mnts[i], f)
		bod, err := os.ReadFile(target)

		Convey("The file is created and that middleware time inserter function works ", t, func() {
			Convey("using a file of ./testdata/j.log", func() {
				Convey("The file is found and the written time matches the correct format", func() {
					So(err, ShouldBeNil)
					So(timeReg.MatchString(string(bod[:len(bod)-2])), ShouldBeTrue)
				})
			})
		})
		os.Remove(target)
	}
}
