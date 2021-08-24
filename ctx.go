package BaiduNetDisk

import (
	"regexp"
)


var (
	fileRegexp *regexp.Regexp
	speCharRegexp *regexp.Regexp
	verifyCsRegexp *regexp.Regexp
)

func init() {
	initRegexp()
}

func initRegexp() {
	var err error
	fileRegexp, err = regexp.Compile(`locals.mset\((.*)\)`)
	if err != nil {
		panic(err)
	}

	speCharRegexp = regexp.MustCompile("\\?|\\||\"|>|<|:|\\*|/|&|#|;|'|\\\\|（|）| |\\+|\\.")
	verifyCsRegexp, err = regexp.Compile("(BDCLND=[^;]+);")
	if err != nil {
		panic(err)
	}
}