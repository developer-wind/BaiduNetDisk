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

	panUrl := "https://pan.baidu.com/s/1Z9tStM7Y5nx59_zIcnzyiw"
	fmt.Println(pu.Size(panUrl, ""))
}