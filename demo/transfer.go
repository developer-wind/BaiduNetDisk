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

	panUrl := "https://pan.baidu.com/s/1iv6js7m0m-qaKW6-GlqXRw"
	pass := "6666"
	path := "/game/101"
	fmt.Println(pu.Transfer(panUrl, path, pass))
}