package domain

type User struct {
	ID            uint64
	OpenID        string
	Mobile        string
	NickName      string
	Email         string
	Avatar        string
	Gender        string
	ThumbsCount   int64
	FansCount     int64
	FollowerCount int64
	Signature     string
	Constellation int8
	Province      string
	City          string
}
