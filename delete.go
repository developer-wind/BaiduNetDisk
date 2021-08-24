package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Delete(fNameList []string) (err error) {
	url := fmt.Sprintf("https://pan.baidu.com/api/filemanager?opera=delete&async=2&onnest=fail&channel=chunlei&web=1&app_id=&bdstoken=%s&clienttype=0", pUser.token)
	fNameListStr, err := json.Marshal(fNameList)
	if err != nil {
		return
	}
	b := strings.NewReader(fmt.Sprint("filelist=", string(fNameListStr)))
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", pUser.cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	c, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ci := new(struct{
		Errno int `json:"errno"`
	})
	err = json.Unmarshal(c, ci)
	if err != nil {
		return
	}
	if ci.Errno != 0 {
		err = errors.New(fmt.Sprintf("删除文件失败，错误码:%d", ci.Errno))
		return
	}
	return
}
