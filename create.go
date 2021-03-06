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
	sync.Once
	sync.Mutex
	list map[string]*sync.Mutex
}

type respCreate struct {
	Exists bool
	DirEmpty bool
}

func (u *PanUser) initCreate() {
	u.create.Do(func() {
		u.create.list = make(map[string]*sync.Mutex)
	})
}

func (u *PanUser) CreatePath(p string) (respSct respCreate, err error) {
	u.create.Lock()
	subMu, exists := u.create.list[p]
	if !exists {
		subMu = &sync.Mutex{}
		u.create.list[p] = subMu
	}
	subMu.Lock()
	defer subMu.Unlock()
	u.create.Unlock()

	parentPath := path.Dir(p)
	fileName := path.Base(p)
	fl, err := u.GetFileList(parentPath)
	if err != nil {
		return
	}

	//文件存在
	info, exists := fl[fileName]
	if exists {
		respSct = respCreate{
			Exists: true,
			DirEmpty: info.DirEmpty == 1,
		}
		return
	}
	respSct = respCreate{
		Exists: false,
		DirEmpty: true,
	}
	url := fmt.Sprintf("https://pan.baidu.com/api/create?a=commit&channel=chunlei&web=1&clienttype=0&bdstoken=%s", u.token)
	b := strings.NewReader(fmt.Sprint("isdir=1&path=", p))
	req, err := http.NewRequest("POST", url, b)
	if err != nil {
		return
	}
	req.Header.Set("Referer", fmt.Sprintf("https://pan.baidu.com/disk/home#/all?vmode=list&path=%s", parentPath))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", u.cookie)
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