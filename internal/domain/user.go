package domain

// User领域对象，暂时可以理解成给前端的响应结果？
type User struct {
	Id       int64
	Email    string
	Password string
}
