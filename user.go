package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type PanUser struct {
	username string

	cookie string
	token  string
	appid  string
	uk string

	quota Quota

	create mu
	list FileList
}

func (u *PanUser) writeCookie(cookie string) error {
	u.cookie = cookie
	return u.verifyCookie()
}

func (u *PanUser) loadQuota() (err error) {
	url := fmt.Sprintf("https://pan.baidu.com/api/quota?checkexpire=1&checkfree=1&channel=chunlei&web=1&app_id=&bdstoken=%s&clienttype=0", u.token)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Referer", fmt.Sprintf("https://pan.baidu.com/disk/home?"))
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
		Total int `json:"total"`
		Free int `json:"free"`
		Used int `json:"used"`
	})
	err = json.Unmarshal(c, ci)
	if err != nil {
		return
	}
	if ci.Errno != 0 {
		err = errors.New(fmt.Sprintf("获取可用空间失败 ，错误码:%d", ci.Errno))
		return
	}
	u.quota = Quota{
		Total: ci.Total,
		Free: ci.Free,
		Used: ci.Used,
	}
	return
}

func ImportCookie(p string) (pu *PanUser, err error) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	c, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	pu = new(PanUser)
	err = pu.writeCookie(string(c))
	if err != nil {
		return
	}
	err = pu.loadQuota()
	if err != nil {
		return
	}
	pu.initCreate()
	pu.initList()
	return
}

func ImportCookies(dir string) (us []*PanUser, err error) {
	ds, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, d := range ds {
		if d.IsDir() {
			usTemp, err1 := ImportCookies(dir + "/" + d.Name())
			if err1 != nil {
				return nil, err1
			}
			us = append(us, usTemp...)
		}
		pu, err2 := ImportCookie(dir + "/" + d.Name())
		if err2 != nil {
			return nil, err2
		}
		us = append(us, pu)
	}
	return
}

func (u *PanUser) verifyCookie() error {
	url := "https://pan.baidu.com/disk/home"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Cookie", u.cookie)
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.Header.Get("Location") == "https://pan.baidu.com/login?redirecturl=http%3A%2F%2Fpan.baidu.com%2Fdisk%2Fhome" {
		return errors.New("BaiduNetDisk cookie is expired")
	}

	j := 5000
	c := make([]byte, 0)
	x := make([]byte, j)
	for j > 0 {
		j, err = resp.Body.Read(x)
		if err != nil && err != io.EOF {
			return err
		}
		c = append(c, x[:j]...)
	}

	fileMatches := fileRegexp.FindSubmatch(c)
	if fileMatches == nil {
		return errors.New("please set the correct cookie, file_struct is not found")
	}
	respStruct := new(struct{
		Token string `json:"bdstoken"`
		Uk string `json:"uk"`
		Username string `json:"username"`
	})
	err = json.Unmarshal(fileMatches[1], respStruct)
	if err != nil {
		return err
	}
	u.username = respStruct.Username
	u.token = respStruct.Token
	u.uk = respStruct.Uk
	return nil
}

func (u *PanUser) Username() string {
	return u.username
}

type Quota struct{
	Total int
	Free int
	Used int
}
func (u *PanUser) Quota() Quota {
	return u.quota
}
