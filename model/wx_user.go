package model

// WxUserReq NewWxUser请求wrapper
type WxUserReq struct {
	NickName string `json:"nick_name" example:"飞燕一号"`
	RealName string `json:"real_name" example:"张三"`
	Sex      string `json:"sex" example:"男"`
	Hometown string `json:"hometown" example:"江苏"`
	Phone    string `json:"phone" example:"133333"`
}

// WxUser 小程序用户信息
type WxUser struct {
	WxUserReq
	ID   int64 `json:"user_id"`
	Role int   `json:"role"`
}
