package domain

type User struct {
	ID       uint64
	OpenID   string
	Mobile   string
	NickName string
	Email    string
	Avatar   string
	Gender   string
}
