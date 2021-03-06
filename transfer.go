package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type File struct {
	size int
	u *PanUser
	sKey string
	Csrf string `json:"csrf"`
	Uk string `json:"uk"`
	Username string `json:"username"`
	Loginstate int `json:"loginstate"`
	VipLevel int `json:"vip_level"`
	SkinName string `json:"skinName"`
	Bdstoken string `json:"bdstoken"`
	Photo string `json:"photo"`
	IsVip int `json:"is_vip"`
	IsSvip int `json:"is_svip"`
	IsEvip int `json:"is_evip"`
	Now time.Time `json:"now"`
	XDUSS string `json:"XDUSS"`
	CurrActivityCode int `json:"curr_activity_code"`
	ShowVipAd int `json:"show_vip_ad"`
	SharePhoto string `json:"share_photo"`
	ShareUk string `json:"share_uk"`
	Shareid int `json:"shareid"`
	HitOgc bool `json:"hit_ogc"`
	ExpiredType int `json:"expiredType"`
	Public int `json:"public"`
	Ctime int `json:"ctime"`
	Description string `json:"description"`
	FollowFlag int `json:"followFlag"`
	AccessListFlag bool `json:"access_list_flag"`
	OwnerVipLevel int `json:"owner_vip_level"`
	OwnerVipType int `json:"owner_vip_type"`
	Linkusername string `json:"linkusername"`
	SharePageType string `json:"share_page_type"`
	TitleImg []string `json:"title_img"`
	FileList []ChildFile `json:"file_list"`
	Errortype int `json:"errortype"`
	Errno int `json:"errno"`
	UfcTime int `json:"ufcTime"`
	Self int `json:"self"`
}

type ChildFile struct {
	AppID string `json:"app_id"`
	Category int `json:"category"`
	DeleteFsID string `json:"delete_fs_id"`
	ExtentInt3 string `json:"extent_int3"`
	ExtentTinyint1 string `json:"extent_tinyint1"`
	ExtentTinyint2 string `json:"extent_tinyint2"`
	ExtentTinyint3 string `json:"extent_tinyint3"`
	ExtentTinyint4 string `json:"extent_tinyint4"`
	FileKey string `json:"file_key"`
	FsID int64 `json:"fs_id"`
	Isdelete string `json:"isdelete"`
	Isdir int `json:"isdir"`
	LocalCtime int `json:"local_ctime"`
	LocalMtime int `json:"local_mtime"`
	Md5 string `json:"md5"`
	OperID string `json:"oper_id"`
	OwnerID string `json:"owner_id"`
	OwnerType string `json:"owner_type"`
	ParentPath string `json:"parent_path"`
	Path string `json:"path"`
	PathMd5 string `json:"path_md5"`
	Privacy string `json:"privacy"`
	RealCategory string `json:"real_category"`
	RootNs int64 `json:"root_ns"`
	ServerAtime string `json:"server_atime"`
	ServerCtime int `json:"server_ctime"`
	ServerFilename string `json:"server_filename"`
	ServerMtime int `json:"server_mtime"`
	Share string `json:"share"`
	Size int `json:"size"`
	Status string `json:"status"`
	TkbindID string `json:"tkbind_id"`
	Videotag string `json:"videotag"`
	Wpfile string `json:"wpfile"`
}

var fileInfoList sync.Map
var verifyList sync.Map

func GetFileInfo(url string) (f *File, err error) {
	r, ok := fileInfoList.Load(url)
	if ok {
		f, _ = r.(*File)
		return
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	//req.Header.Set("Cookie", u.cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	i,j := 500, 500
	c := make([]byte, 0)
	x := make([]byte, j)
	for i == j {
		j, err = resp.Body.Read(x)
		if err != nil && err != io.EOF {
			return
		}
		c = append(c, x[:j]...)
	}

	if fileDelRegexp.Match(c) {
		err = ErrDel
		return
	}
	fileMatches := fileRegexp.FindSubmatch(c)
	if fileMatches == nil {
		err = errors.New(fmt.Sprintln(url, ": file_struct is not found"))
		return
	}

	f = new(File)
	//f.u = u
	err = json.Unmarshal(fileMatches[1], f)
	if err != nil {
		return
	}
	fileInfoList.Store(url, f)
	return
}

func (f *File) Size() (size int, err error) {
	if f.size > 0 {
		size = f.size
		return
	}
	var getListSize func(fs []FileInfo) (size int, err error)
	getListSize = func(fs []FileInfo) (size int, err error) {
		for _, info := range fs {
			if info.Isdir == 1 {
				if info.DirEmpty == 1 {
					continue
				}
				fsTemp, err1 := f.list(info.Path)
				if err1 != nil {
					err = err1
					return
				}
				sizeTemp, err2 := getListSize(fsTemp)
				if err2 != nil {
					err = err2
					return
				}
				size += sizeTemp
				continue
			}
			size += info.Size
		}
		return
	}

	for _, file := range f.FileList {
		if file.Isdir == 1 {
			fs, err3 := f.list(file.Path)
			if err3 != nil {
				err = err3
				return
			}
			sizeTemp, err4 := getListSize(fs)
			if err4 != nil {
				err = err4
				return
			}
			size += sizeTemp
			continue
		}
		size += file.Size
	}
	f.size = size
	return
}

func (f *File) Verify(url, pass string) (err error) {
	if pass == "" {
		return
	}
	_, ok := verifyList.Load(url+pass)
	if ok {
		return
	}

	sUrl := strings.Replace(url, "https://pan.baidu.com/s/1", "", -1)
	url1 := fmt.Sprintf("https://pan.baidu.com/share/verify?surl=%s&channel=chunlei&web=1&bdstoken=%s&clienttype=0", sUrl, f.Bdstoken)
	b := strings.NewReader(fmt.Sprint("pwd=", pass))
	req, err := http.NewRequest("POST", url1, b)
	if err != nil {
		return
	}
	req.Header.Set("Referer", fmt.Sprintf("https://pan.baidu.com/share/init?surl=%s", sUrl))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", f.u.cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	c, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	respObj := new(struct {
		Errno int `json:"errno"`
		ErrMsg string `json:"err_msg"`
		Randsk string `json:"randsk"`
	})
	err = json.Unmarshal(c, respObj)
	if err != nil {
		return
	}
	if respObj.Errno != 0 {
		err = errors.New(fmt.Sprintf("%s ??????????????????????????????:%d", url, respObj.Errno))
		return
	}
	f.sKey = respObj.Randsk
	tokenMatches := verifyCsRegexp.FindStringSubmatch(resp.Header.Get("Set-Cookie"))
	if tokenMatches == nil {
		err = errors.New(fmt.Sprintf("%s ????????????-??????Cookie??????????????????:%d", url, respObj.Errno))
		return
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	cs := fmt.Sprintf("%s;%s", tokenMatches[1], f.u.cookie)
	req.Header.Set("Cookie", cs)
	resp, err = (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	i,j := 500, 500
	c = make([]byte, 0)
	x := make([]byte, j)
	for i == j {
		j, err = resp.Body.Read(x)
		if err != nil && err != io.EOF {
			return
		}
		c = append(c, x[:j]...)
	}

	fileMatches := fileRegexp.FindSubmatch(c)
	if fileMatches == nil {
		err = errors.New(fmt.Sprintln(url, ": file_struct is not found"))
		return
	}
	f1 := new(File)
	err = json.Unmarshal(fileMatches[1], f1)
	if err != nil {
		return
	}
	f.FileList = f1.FileList
	f.Shareid = f1.Shareid
	f.Bdstoken = f1.Bdstoken
	verifyList.Store(url+pass, struct {}{})
	return
}

func Size(url, pass string) (size int, err error) {
	f, err := GetFileInfo(url)
	if err != nil {
		return
	}

	if pass == "" && len(f.FileList) == 0 {
		err = errors.New("?????????????????????????????????????????????")
		return
	}
	err = f.Verify(url, pass)
	if err != nil {
		return
	}

	size, err = f.Size()
	if err != nil {
		return
	}
	return
}

var InsufficientSpaceError = errors.New("??????????????????,????????????")
func (u *PanUser) Transfer(url, path, pass string) error {
	//????????????????????????
	//??????????????????????????????????????????????????????????????????????????????
	f, err := GetFileInfo(url)
	if err != nil {
		return err
	}
	f.u = u

	if pass == "" && len(f.FileList) == 0 {
		return errors.New("?????????????????????????????????????????????")
	}

	respCreate, err := u.CreatePath(path)
	if err != nil {
		return err
	}
	if respCreate.Exists && !respCreate.DirEmpty {
		return os.ErrExist
	}

	err = f.Verify(url, pass)
	if err != nil {
		return err
	}

	fSidList := make([]int64, 0)
	for _, list := range f.FileList {
		fSidList = append(fSidList, list.FsID)
	}
	fSidListBts, err := json.Marshal(fSidList)
	if err != nil {
		return err
	}

	b := strings.NewReader(fmt.Sprint("fsidlist=", string(fSidListBts), "&path=", path))
	url2 := fmt.Sprintf("https://pan.baidu.com/share/transfer?shareid=%d&from=%s&ondup=newcopy&async=1&channel=chunlei&web=1&app_id=%s&bdstoken=%s&clienttype=0&sekey=%s", f.Shareid, f.ShareUk, u.appid, u.token, f.sKey)
	req, err := http.NewRequest("POST", url2, b)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", url)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", u.cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	c, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respObj := new(struct {
		Errno int `json:"errno"`
	})
	err = json.Unmarshal(c, respObj)
	if err != nil {
		return err
	}
	if respObj.Errno == 12 {
		err = InsufficientSpaceError
		return err
	}
	if respObj.Errno != 0 {
		err = errors.New(fmt.Sprintf("?????????:%d %s",respObj.Errno, path))
		return err
	}
	return nil
}