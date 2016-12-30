package duobb

import (
	"github.com/reechou/holmes"
	"github.com/reechou/tao"
)

func (self *DuobbProcess) PushMsg(method string, user, msg []byte, conn tao.Connection) {
	rsp := &DuobbMsg{
		UserName: user,
		Method: []byte(method),
		Msg: msg,
	}
	switch method {
	case DUOBB_ACCESS_LOGOUT_KICKOFF:
	default:
		holmes.Error("push msg error: cannot found method[%s]", method)
		return
	}
	err := conn.Write(rsp)
	if err != nil {
		holmes.Error("conn[%s] write response msg[%v] error: %v", conn.GetName(), rsp, err)
	}
}

func (self *DuobbProcess) createLogoutKickoffMsg() []byte {
	return nil
}
