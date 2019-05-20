package model

// WxUserInfo 非重要个人信息
type WxUserInfo struct {
	NickName string `json:"nick_name" example:"飞燕一号"`
	RealName string `json:"real_name" example:"张三"`
	Sex      string `json:"sex" example:"男"`
	Hometown string `json:"hometown" example:"江苏"`
	Phone    string `json:"phone" example:"133333"`
}

// WxUserImportantInfo 重要个人信息
type WxUserImportantInfo struct {
	OpenID     string `json:"open_id" example:"jajsjasja11233"`
	SessionKey string `json:"session_key" example:"jajsjasja11233"`
}

// WxUser 用户信息数据库字段
type WxUser struct {
	ID int64 `json:"user_id"`
	WxUserInfo
	WxUserImportantInfo
	Role       int
	UpdateTime int64
}
