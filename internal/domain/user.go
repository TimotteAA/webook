package domain

// User领域对象，暂时可以理解成给前端的响应结果？
type User struct {
	Id          int64  `json:"id"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"-"`
	Description string `json:"description,omitempty"`
	NickName    string `json:"nickname,omitempty"`
	BirthDay    string `json:"birthday,omitempty"`
}
