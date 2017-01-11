package duobb

const (
	LOGIN = iota
	LOGOUT
	LOGOUT_KICKOFF
)

const (
	MAX_LOGOUT_KICKOFF = 3
)

type Session struct {
	User        string
	Status      int
	CheckLogout int // 检测登出次数
	AppId       int
	Version     string
}
