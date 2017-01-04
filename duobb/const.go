package duobb

const (
	CONN_NAME_DELIMITER = "_$^%*$_"
)

const (
	BACKEND_GOJSONRPC   = "gojsonrpc"
	BACKEND_NORMAL_HTTP = "normalhttp"
)

const (
	DUOBB_ACCESS_LOGIN          = "DuobbAccountService.Login"
	DUOBB_ACCESS_LOGOUT         = "DuobbAccountService.Logout"
	DUOBB_ACCESS_HEARTBEAT      = "DuobbAccountService.Heartbeat"
	DUOBB_ACCESS_GETALLDATA     = "DuobbAccountService.GetAllDuobbData"
	DUOBB_ACCESS_LOGOUT_KICKOFF = "DuobbAccountService.LogoutKickOff"
)

const (
	APPID_DINGJIAN_TAOKE = iota
	APPID_DINGJIAN_QUNGUAN
)

var apps []int = []int{APPID_DINGJIAN_TAOKE, APPID_DINGJIAN_QUNGUAN}
