package BaiduNetDisk

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type PanUser struct {
	cookie string
	token  string
	appid  string
	uk string
}

var pUser = new(PanUser)

func WriteCookie(cookie string) error {
	pUser.cookie = cookie
	return pUser.verifyCookie()
}

func ImportCookie(p string) error {
	f, e := os.Open(p)
	if e != nil {
		return e
	}
	c, er := ioutil.ReadAll(f)
	if er != nil {
		return er
	}
	return WriteCookie(string(c))
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
		return errors.New("file_struct is not found")
	}
	respStruct := new(struct{
		Token string `json:"bdstoken"`
		Uk string `json:"uk"`
	})
	err = json.Unmarshal(fileMatches[1], respStruct)
	if err != nil {
		return err
	}
	u.token = respStruct.Token
	u.uk = respStruct.Uk
	return nil
}
