package ext

const (
	DUOBB_MANAGER_RESPONSE_OK = 1000
)

const (
	WX_SEND_MSG_CHAN_LEN = 1024
)

type DuobbManagerReq struct {
	Uri string
	Req interface{}
}

type DuobbManagerSendMsgReq struct {
	UserName string `json:"userName"`
	Msg      string `json:"meg"`
}

type DuobbManagerResponse struct {
	Code int64       `json:"state"`
	Msg  string      `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}
