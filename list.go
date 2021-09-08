package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type FileList struct {
	sync.Once
	sync.Mutex
	List map[string]*struct{
		sync.Mutex
		List map[string]FileInfo
	}
}

type FileListResp struct {
	Errno int `json:"errno"`
	List []FileInfo `json:"list"`
}
type FileInfo struct {
	ServerFilename string `json:"server_filename"`
	Isdir int `json:"isdir"`
	DirEmpty int `json:"dir_empty"`
	Path string `json:"path"`
	OwnerID int `json:"owner_id"`
	Size int `json:"size"`
	Share int `json:"share"`
	Empty int `json:"empty"`
	FsID int64 `json:"fs_id"`
}

const ParentPath = "/"

func (u *PanUser) initList() {
	u.list.Do(func() {
		u.list.List = make(map[string]*struct{
			sync.Mutex
			List map[string]FileInfo
		})
	})
}

func (u *PanUser) GetFileList(path string) (map[string]FileInfo, error) {
	if path == "" {
		path = ParentPath
	}

	//构建子结构，降低锁竞争
	u.list.Lock()
	fl, exists := u.list.List[path]
	if !exists {
		fl = new(struct{
			sync.Mutex
			List map[string]FileInfo
		})
		fl.List = make(map[string]FileInfo)
		u.list.List[path] = fl
	}
	u.list.Unlock()

	//子结构加锁
	fl.Lock()
	defer fl.Unlock()
	if len(fl.List) > 0 {
		return fl.List, nil
	}

	var getListFun func(page int) (err error)
	getListFun = func(page int) (err error) {
		limit := 1000
		start := (page - 1) * limit
		url := fmt.Sprintf("https://pan.baidu.com/api/list?start=%d&limit=%d&channel=chunlei&web=1&clienttype=0&dir=%s", start, limit, path)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.Header.Set("Referer", fmt.Sprintf("https://pan.baidu.com/disk/home?#/all?vmode=list&path=%s", path))
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
		flResp := new(FileListResp)
		err = json.Unmarshal(c, flResp)
		if err != nil {
			return
		}
		if flResp.Errno != 0 {
			err = errors.New(fmt.Sprintf("获取文件列表失败，错误码：%d", flResp.Errno))
			return
		}
		for _, list := range flResp.List {
			fl.List[list.ServerFilename] = list
		}
		if len(flResp.List) == limit {
			page ++
			err = getListFun(page)
			if err != nil {
				return
			}
		}
		return
	}

	err := getListFun(1)
	if err != nil {
		return nil, err
	}
	return fl.List, nil
}

func (f *File) list(path string) (fs []FileInfo, err error) {
	var getListFun func(page int) (err error)
	getListFun = func(page int) (err error) {
		limit := 1000
		url := fmt.Sprintf("https://pan.baidu.com/share/list?uk=%s&shareid=%d&order=other&desc=1&showempty=0&web=1&page=%d&num=%d&dir=%s&channel=chunlei&web=1&app_id=&bdstoken=%s&clienttype=0",
			f.ShareUk,f.Shareid,page,limit,path,f.Bdstoken,
		)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return
		}
		req.Header.Set("Referer", fmt.Sprintf("https://pan.baidu.com/disk/home?#/all?vmode=list"))
		//req.Header.Set("Cookie", f.u.cookie)
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		c, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		flResp := new(FileListResp)
		err = json.Unmarshal(c, flResp)
		if err != nil {
			return
		}
		if flResp.Errno != 0 {
			err = errors.New(fmt.Sprintf("获取分享文件列表失败，错误码：%d", flResp.Errno))
			return
		}
		fs = append(fs, flResp.List...)
		if len(flResp.List) == limit {
			page ++
			err = getListFun(page)
			if err != nil {
				return
			}
		}
		return
	}

	err = getListFun(1)
	if err != nil {
		return nil, err
	}
	return
}