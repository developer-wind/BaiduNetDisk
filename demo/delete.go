package main

import (
	"fmt"
	"github.com/developer-wind/BaiduNetDisk"
)

func main() {
	err := BaiduNetDisk.ImportCookie("/Users/tt/Downloads/game_down/pan.cookies")
	if err != nil {
		panic(err)
	}

	fmt.Println(BaiduNetDisk.Delete([]string{
		"/game/101xx",
		"/game/xxx",
	}))
}


