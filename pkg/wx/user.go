package wx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Session struct {
	OpenID     string `json:"openid"`
	Errcode    int    `json:"errcode"`
	SessionKey string `json:"session_key"`
	ErrMsg     string `json:"errmsg"`
	Unionid    string `json:"unionid"`
}

func Code2Session(appID, secret, code string) (*Session, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appID, secret, code)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r Session
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}
	if r.Errcode != 0 {
		return nil, fmt.Errorf("auth code error: %s", r.ErrMsg)
	}
	return &r, nil
}

type UserInfo struct {
	OpenID    string `json:"openid" xorm:"openid"`
	NickName  string `json:"nickName" xorm:"nickName"`
	Gender    string `json:"gender" xorm:"gender"`
	Province  string `json:"province" xorm:"province"`
	City      string `json:"city" xorm:"city"`
	Country   string `json:"country" xorm:"country"`
	AvatarUrl string `json:"avatarUrl" xorm:"avatarUrl"`
}

func BizDataCrypt(appID, sessionKey, iv, encryptedData string) (*UserInfo, error) {
	pc := WxBizDataCrypt{AppID: appID, SessionKey: sessionKey}
	result, err := pc.Decrypt(encryptedData, iv, true)
	if err != nil {
		return nil, err
	}

	var u UserInfo
	err = json.Unmarshal([]byte(result.(string)), &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
