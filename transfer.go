package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type File struct {
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

func GetFileInfo(url string) (f *File, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Cookie", pUser.cookie)
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

	fileMatches := fileRegexp.FindSubmatch(c)
	if fileMatches == nil {
		err = errors.New(fmt.Sprintln(url, ": file_struct is not found"))
		return
	}

	f = new(File)
	err = json.Unmarshal(fileMatches[1], f)
	if err != nil {
		return
	}
	return
}

func (f *File) Verify(url, pass string) (err error) {
	if pass == "" {
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
		err = errors.New(fmt.Sprintf("%s 密码验证失败，错误码:%d", url, respObj.Errno))
		return
	}
	f.sKey = respObj.Randsk
	tokenMatches := verifyCsRegexp.FindStringSubmatch(resp.Header.Get("Set-Cookie"))
	if tokenMatches == nil {
		err = errors.New(fmt.Sprintf("%s 密码验证-获取Cookie失败，错误码:%d", url, respObj.Errno))
		return
	}

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	cs := fmt.Sprintf("%s;%s", tokenMatches[1], pUser.cookie)
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
	return
}

func Transfer(url, path, pass string) error {
	//获取文件相关参数
	//需要提取码的文件，部分参数获取不到，需在验证环节获取
	f, err := GetFileInfo(url)
	if err != nil {
		return err
	}

	respCreate, err := CreatePath(path)
	if err != nil {
		return err
	}
	if respCreate.exists && !respCreate.dirEmpty {
		return os.ErrExist
	}

	err = f.Verify(url, pass)
	if err != nil {
		return err
	}
	//需要提取码的文件，需验证过后FileList才是真实数据
	if len(f.FileList) == 0 {
		return ErrDel
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
	url2 := fmt.Sprintf("https://pan.baidu.com/share/transfer?shareid=%d&from=%s&ondup=newcopy&async=1&channel=chunlei&web=1&app_id=%s&bdstoken=%s&clienttype=0&sekey=%s", f.Shareid, f.ShareUk, pUser.appid, pUser.token, f.sKey)
	req, err := http.NewRequest("POST", url2, b)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", url)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", pUser.cookie)
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
	if respObj.Errno != 0 {
		err = errors.New(fmt.Sprintf("错误码:%d %s",respObj.Errno, path))
		return err
	}
	return nil
}