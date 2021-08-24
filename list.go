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
	mu sync.Mutex
	List map[string]*struct{
		mu sync.Mutex
		List map[string]FileInfo
	}
}

type FileListResp struct {
	Errno int `json:"errno"`
	GUIDInfo string `json:"guid_info"`
	List []FileInfo `json:"list"`
	RequestID int64 `json:"request_id"`
	GUID int `json:"guid"`
}
type FileInfo struct {
	TkbindID int `json:"tkbind_id"`
	ServerFilename string `json:"server_filename"`
	OwnerType int `json:"owner_type"`
	Category int `json:"category"`
	RealCategory string `json:"real_category"`
	Isdir int `json:"isdir"`
	DirEmpty int `json:"dir_empty"`
	Path string `json:"path"`
	Wpfile int `json:"wpfile"`
	OperID int64 `json:"oper_id"`
	ServerCtime int `json:"server_ctime"`
	OwnerID int `json:"owner_id"`
	LocalMtime int `json:"local_mtime"`
	Size int `json:"size"`
	Unlist int `json:"unlist"`
	Share int `json:"share"`
	ServerMtime int `json:"server_mtime"`
	Pl int `json:"pl"`
	LocalCtime int `json:"local_ctime"`
	ServerAtime int `json:"server_atime"`
	Empty int `json:"empty"`
	FsID int64 `json:"fs_id"`
}

const ParentPath = "/"

var (
	FList = new(FileList)
)
func init() {
	FList.List = make(map[string]*struct{
		mu sync.Mutex
		List map[string]FileInfo
	})
}

func GetFileList(path string) (map[string]FileInfo, error) {
	if path == "" {
		path = ParentPath
	}

	//构建子结构，降低锁竞争
	FList.mu.Lock()
	fl, exists := FList.List[path]
	if !exists {
		fl = new(struct{
			mu sync.Mutex
			List map[string]FileInfo
		})
		fl.List = make(map[string]FileInfo)
		FList.List[path] = fl
	}
	FList.mu.Unlock()

	//子结构加锁
	fl.mu.Lock()
	defer fl.mu.Unlock()
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
		req.Header.Set("Cookie", pUser.cookie)
		resp, err := (&http.Client{}).Do(req)
		if err != nil {
			return
		}
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