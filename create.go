package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
)

type mu struct {
	mu sync.Mutex
	list map[string]*sync.Mutex
}

type respCreate struct {
	exists bool
	dirEmpty bool
}

var (
	m = new(mu)
)

func init() {
	m.list = make(map[string]*sync.Mutex)
}

func CreatePath(p string) (respSct respCreate, err error) {
	m.mu.Lock()
	subMu, exists := m.list[p]
	if !exists {
		subMu = &sync.Mutex{}
		m.list[p] = subMu
	}
	subMu.Lock()
	defer subMu.Unlock()
	m.mu.Unlock()

	parentPath := path.Dir(p)
	fileName := path.Base(p)
	fl, err := GetFileList(parentPath)
	if err != nil {
		return
	}

	//文件存在
	info, exists := fl[fileName]
	if exists {
		respSct = respCreate{
			exists: true,
			dirEmpty: info.DirEmpty == 1,
		}
		return
	}
	respSct = respCreate{
		exists: false,
		dirEmpty: true,
	}
	url := fmt.Sprintf("https://pan.baidu.com/api/create?a=commit&channel=chunlei&web=1&clienttype=0&bdstoken=%s", pUser.token)
	b := strings.NewReader(fmt.Sprint("isdir=1&path=", p))
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return
	}
	req.Header.Set("Referer", fmt.Sprintf("https://pan.baidu.com/disk/home#/all?vmode=list&path=%s", parentPath))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", pUser.cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
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
		err = errors.New(fmt.Sprintf("%s 创建文件失败，错误码:%d", p, ci.Errno))
		return
	}
	return
}