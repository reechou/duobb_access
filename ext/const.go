package ext

const (
	DUOBB_MANAGER_TYPE_LOGIN = iota
	DUOBB_MANAGER_TYPE_LOGOUT
	DUOBB_MANAGER_TYPE_HEALTH_ERROR
)

const (
	DUOBB_MANAGER_WECHAT_MSG_URI = "/index.php?r=report/sendtemp"
)

const (
	LOGIN_MSG    = "顶尖淘客 [%s] 登录成功!"
	LOGOUT_MSG   = "顶尖淘客 [%s] 已正常登出!"
	HEALTH_ERROR = "健康检查出错, 顶尖淘客 [%s] 已异常登出, 请检查!"
)
