package duobb

import (
	"fmt"
	"strings"
	"time"

	"github.com/reechou/duobb_proto"
	"github.com/reechou/holmes"
	"github.com/reechou/tao"
)

// return false: login or logout request or not login
// return true : other request of login true
func (self *DuobbProcess) checkMsg(msg *DuobbMsg, secretKey []byte, decodeMsg []byte, conn tao.Connection) bool {
	requestUser := string(msg.UserName)
	c := self.connMap[requestUser]
	rsp := &DuobbMsg{
		UserName: msg.UserName,
		Method:   msg.Method,
	}
	resultResponse := &duobb_proto.Response{
		Code: duobb_proto.DUOBB_RSP_SUCCESS,
	}
	ifLogout := false
	switch string(msg.Method) {
	case DUOBB_ACCESS_LOGIN:
		user, err := self.checkReqMsgUser(decodeMsg)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGIN_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGIN_ERROR
		} else {
			if user != requestUser {
				holmes.Error("user[%s] not equal requestUser[%s]", user, requestUser)
				resultResponse.Code = duobb_proto.DUOBB_MSG_LOGIN_ERROR
				resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGIN_ERROR
			} else {
				if c != nil {
					if c.GetName() != conn.GetName() {
						// maybe write [kick off line] msg
						holmes.Info("client[%s] kick off line, for relogin", c.GetName())
						c.Close()

						self.connMutex.Lock()
						conn.SetName(requestUser + CONN_NAME_DELIMITER + conn.GetName())
						holmes.Info("user[%s] login success with conn name[%s]", requestUser, conn.GetName())
						self.connMap[requestUser] = conn
						self.connMutex.Unlock()
					}
				} else {
					self.connMutex.Lock()
					conn.SetName(requestUser + CONN_NAME_DELIMITER + conn.GetName())
					self.connMap[requestUser] = conn
					self.connMutex.Unlock()
					holmes.Info("user[%s] login success with conn name[%s]", requestUser, conn.GetName())
				}
			}
		}
	case DUOBB_ACCESS_LOGOUT:
		user, err := self.checkReqMsgUser(decodeMsg)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
		} else {
			if user != requestUser {
				holmes.Error("user[%s] not equal requestUser[%s]", user, requestUser)
				resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
				resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
			} else {
				if c != nil {
					if c.GetName() == conn.GetName() {
						ifLogout = true
						holmes.Info("start to logout: %s", c.GetName())
					} else {
						holmes.Error("user[%s] conn.name[%s] not equal connmap.name[%s]", requestUser, conn.GetName(), c.GetName())
						resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
						resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
					}
				} else {
					holmes.Error("user[%s] has none this user in connmap", requestUser)
					resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
					resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
				}
			}
		}
	case DUOBB_ACCESS_HEARTBEAT:
		holmes.Debug("conn[%s] in heartbeat", conn.GetName())
		heartbeatData, err := self.checkHeartbeat(decodeMsg)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_HEARTBEAT_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_HEARTBEAT_ERROR
		} else {
			resultResponse.Data = heartbeatData
		}
	case DUOBB_ACCESS_GETALLDATA:
		return true
	default:
		if c == nil {
			holmes.Error("user[%s] has no login.", string(msg.UserName))
			// clear conn
			conn.Close()
			return false
		}
		if c.GetName() != conn.GetName() {
			holmes.Error("user[%s] now conn[%s] not equal connmap conn[%s], close conn.", requestUser, conn.GetName(), c.GetName())
			// clear conn
			conn.Close()
			return false
		}
		return true
	}
	resultMsg, err := JsonEncode(resultResponse)
	if err != nil {
		holmes.Error("json encode[%v] error: %v", resultResponse, err)
		return false
	}
	resultMsg, err = self.encodeMsg(secretKey, resultMsg)
	if err != nil {
		holmes.Error("user[%s] response msg encode error: %v", string(msg.UserName), err)
		return false
	}
	rsp.Msg = resultMsg
	err = conn.Write(rsp)
	if err != nil {
		holmes.Error("conn[%s] write response msg[%v] error: %v", conn.GetName(), rsp, err)
	}
	//if resultResponse.Code != duobb_proto.DUOBB_RSP_SUCCESS {
	//	// clear conn
	//	conn.Close()
	//}
	if ifLogout {
		// clear conn
		c.Close()
		holmes.Info("user[%s] conn.name[%s] logout success", requestUser, c.GetName())
	}

	return false
}

func (self *DuobbProcess) checkReqMsgUser(decodeMsg []byte) (string, error) {
	request, err := JsonDecode(decodeMsg)
	if err != nil {
		holmes.Error("json decode[%s] error: %v", string(decodeMsg), err)
		return "", err
	}
	reqMap := request.(map[string]interface{})
	if reqMap == nil {
		holmes.Error("request: %v translate to map error", request)
		return "", fmt.Errorf("request: %v translate to map error", request)
	}
	user := reqMap["user"]
	if user == nil {
		holmes.Error("reqmap: %v has no user field", reqMap)
		return "", fmt.Errorf("reqmap: %v has no user field", reqMap)
	}
	userStr, ok := user.(string)
	if !ok {
		holmes.Error("user: %v translate to string error", user)
		return "", fmt.Errorf("user: %v translate to string error", user)
	}
	return userStr, nil
}

func (self *DuobbProcess) OnConnectCallback(conn tao.Connection) bool {
	holmes.Info("[%s] on connect", conn.GetName())
	return true
}

func (self *DuobbProcess) OnErrorCallback() {
	holmes.Info("on Error")
}

func (self *DuobbProcess) OnCloseCallback(conn tao.Connection) {
	connName := conn.GetName()
	holmes.Info("conn[%s] on close", connName)
	names := strings.Split(conn.GetName(), CONN_NAME_DELIMITER)
	if len(names) != 2 {
		holmes.Error("[%s] name error.", connName)
		return
	}
	self.connMutex.Lock()
	defer self.connMutex.Unlock()
	c := self.connMap[names[0]]
	if c == nil {
		holmes.Error("[%s] connmap has no this connection.", connName)
		return
	}
	delete(self.connMap, names[0])
	holmes.Info("conn[%s] on close clear success.", connName)
}

func (self *DuobbProcess) OnScheduleCallback(now time.Time, data interface{}) {
	conn := data.(tao.Connection)
	lastTimestamp := conn.GetHeartBeat()
	if now.UnixNano()-lastTimestamp > int64(self.cfg.ServerConfig.Timeout)*1e9 {
		holmes.Warn("Client %d %s timeout, close it\n", conn.GetNetId(), conn.GetName())
		conn.Close()
	}
}
