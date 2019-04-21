package model

// WxUserWrapper NewWxUser请求wrapper
type WxUserWrapper struct {
	NickName string `json:"nick_name" example:"飞燕一号"`
	RealName string `json:"real_name" example:"张三"`
	Sex      string `json:"sex" example:"男"`
	Hometown string `json:"hometown" example:"江苏"`
	Phone    string `json:"phone" example:"133333"`
}

// WxUser 小程序用户信息
type WxUser struct {
	WxUserWrapper
	ID   int64 `json:"user_id"`
	Role int   `json:"role"`
}
