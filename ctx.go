package BaiduNetDisk

import (
	"errors"
	"regexp"
)

var (
	fileRegexp *regexp.Regexp
	fileDelRegexp *regexp.Regexp
	verifyCsRegexp *regexp.Regexp
	ErrDel = errors.New("分享的文件已被删除")
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

	verifyCsRegexp, err = regexp.Compile(`(BDCLND=[^;]+);`)
	if err != nil {
		panic(err)
	}

	fileDelRegexp, err = regexp.Compile(`分享的文件已经被删除`)
	if err != nil {
		panic(err)
	}
}