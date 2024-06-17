package user

type User struct {
	Id       int64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
}
