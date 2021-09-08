package main

import (
	"fmt"
	"github.com/developer-wind/BaiduNetDisk"
)

func main() {
	pu, err := BaiduNetDisk.ImportCookie("/Users/tt/Downloads/game_down/pan.cookies")
	if err != nil {
		panic(err)
	}

	fmt.Println(pu.GetFileList("/game"))
}


