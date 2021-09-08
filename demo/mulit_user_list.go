package main

import (
	"fmt"
	"github.com/developer-wind/BaiduNetDisk"
)

func main() {
	us, err := BaiduNetDisk.ImportCookies("/Users/tt/Downloads/pan_cookies")
	if err != nil {
		panic(err)
	}

	for _, u := range us {
		fmt.Print(u.Username(), " ")
		fmt.Print(u.Quota(), " ")
		fmt.Println(u.GetFileList("/"))
	}
}
