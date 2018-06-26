package model

import (
	"encoding/json"
)

//人物的个人信息数据
type Profile struct {
	//具体的人物信息
	Name       string
	Gender     string
	Age        int
	Height     int
	Weight     int
	Income     string
	Marriage   string
	Education  string
	Occupation string
	Hokou      string
	Xinzuo     string
	House      string
	Car        string
}

func FromJsonObj(o interface{}) (Profile, error) {
	var profile Profile
	//	将数据编码成json字符串
	s, err := json.Marshal(o)
	if err != nil {
		return profile, err
	}

	//将json字符串解码成Profile
	err = json.Unmarshal(s, &profile)
	return profile, err
}
