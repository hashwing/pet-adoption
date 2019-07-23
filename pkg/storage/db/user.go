package db

import "errors"

type User struct {
	ID        string `json:"uuid" xorm:"'uuid'"`
	OpenID    string `json:"openid" xorm:"openid"`
	NickName  string `json:"nick_name" xorm:"nick_name"`
	Gender    string `json:"gender" xorm:"gender"`
	Province  string `json:"province" xorm:"province"`
	City      string `json:"city" xorm:"city"`
	Country   string `json:"country" xorm:"country"`
	AvatarUrl string `json:"avatar_url" xorm:"avatar_url"`
}

func AddUser(user User) error {
	_, err := MysqlDB.Insert(user)
	return err
}

func UpdateUser(user User) error {
	_, err := MysqlDB.Where("uuid=?", user.ID).Update(user)
	return err
}

func UserExistByOpenID(openid string) (*User, bool, error) {
	var u User
	isExist, err := MysqlDB.Where("openid=?", openid).Get(&u)
	return &u, isExist, err
}

func GetUser(uuid string) (*User, error) {
	var u User
	isExist, err := MysqlDB.Where("uuid=?", uuid).Get(&u)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.New("user not found")
	}
	return &u, nil
}
