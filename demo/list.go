package main

import (
	"fmt"
	"github.com/developer-wind/BaiduNetDisk"
)

func main() {
	pu, err := BaiduNetDisk.ImportCookie("/Users/tt/Downloads/pan_cookies/qq2215219591")
	if err != nil {
		panic(err)
	}

	fmt.Println(pu.GetFileList("/game"))
}


