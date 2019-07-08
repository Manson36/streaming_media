package defs

//requests
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd string 		`json:"pwd"`
}

type NewComment struct {
	AuthorId int  	 `json:"author_id"`
	Content string	 `json:"content"`
}

type NewVideo struct {
	AuthorId int	 `json:"author_id"`
	Name string 	 `json:"name"`
}

//response
type SignedUp struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}

type UserSession struct {
	Username string `json:"user_name"`
	SessionId string `json:"session_id"`
}

type UserInfo struct {
	Id int `json:"id"`
}

type SignedIn struct {
	Success bool `json:"success"`
	SessionId string `json:"session_id"`
}

type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

type Comments struct {
	Comments []*comment `json:"comments"`
}

//Data model
type User struct {
	Id int
	LoginName string
	Pwd string
}

type comment struct {
	Id string `json:"id"`
	VideoId string `json:"video_id"`
	Author string `json:"author"`
	Content string `json:"content"`
}

type SimpleSession struct {
	Username string //login name
	TTL int64 //用来检查用户是否过期
}

type VideoInfo struct {
	Id string
	AuthorId int
	Name string
	DisplayCtime string
}
