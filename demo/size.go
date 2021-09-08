package main

import (
	"fmt"
	"github.com/developer-wind/BaiduNetDisk"
)

func main() {
	panUrl := "https://pan.baidu.com/s/1Z9tStM7Y5nx59_zIcnzyiw"
	fmt.Println(BaiduNetDisk.Size(panUrl, ""))
}