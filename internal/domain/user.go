package domain

// User领域对象，暂时可以理解成给前端的响应结果？
type User struct {
	Id            int64  `json:"id"`
	Email         string `json:"email,omitempty"`
	Password      string `json:"-"`
	Description   string `json:"description,omitempty"`
	NickName      string `json:"nickname,omitempty"`
	BirthDay      int64  `json:"birthday,omitempty"`
	Phone         string `json:"phone"`
	CreatetAt     int64  `json:"createdAt"`
	WeChatOpenId  string `json:"weChatOpenId"`
	WeChatUnionId string `json:"weChatUnionId"`
}
