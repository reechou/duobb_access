package duobb

import (
	"fmt"
	"time"

	"github.com/reechou/duobb_access/ext"
	"github.com/reechou/duobb_proto"
	"github.com/reechou/holmes"
	"github.com/reechou/tao"
)

// return false: login or logout request or not login
// return true : other request of login true
func (self *DuobbProcess) checkMsg(msg *DuobbMsg, secretKey []byte, decodeMsg []byte, conn tao.Connection) bool {
	requestUser := string(msg.UserName)
	appid := int(msg.AppId)
	rsp := &DuobbMsg{
		AppId:    msg.AppId,
		UserName: msg.UserName,
		Method:   msg.Method,
	}
	connMap := self.connMap[appid]
	if connMap == nil {
		holmes.Error("cannot found this appid[%d]", appid)
		return false
	}
	c := connMap[requestUser]
	resultResponse := &duobb_proto.Response{
		Code: duobb_proto.DUOBB_RSP_SUCCESS,
	}
	ifLogout := false
	switch string(msg.Method) {
	case DUOBB_ACCESS_LOGIN:
		reqMap, err := self.parseMsg(decodeMsg)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGIN_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGIN_ERROR
			break
		}
		user, err := self.checkReqMsgUser(reqMap)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGIN_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGIN_ERROR
			break
		}
		if user != requestUser {
			holmes.Error("user[%s] not equal requestUser[%s]", user, requestUser)
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGIN_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGIN_ERROR
			break
		}
		version, err := self.checkReqMsgVersion(reqMap)
		if err != nil {
			version = "0.0.0"
		}
		if c != nil {
			if c.GetName() != conn.GetName() {
				holmes.Info("client[%s] start to kick off line, for relogin", c.GetName())
				extra := c.GetExtraData()
				if extra == nil {
					holmes.Error("get extra data error == nil")
					c.Close()
				} else {
					session := extra.(Session)
					session.Status = LOGOUT_KICKOFF
					c.SetExtraData(session)
					// send kick off msg to client
					self.PushMsg(DUOBB_ACCESS_LOGOUT_KICKOFF, msg.UserName, nil, c)
				}

				self.connMutex.Lock()
				//conn.SetName(requestUser + CONN_NAME_DELIMITER + conn.GetName())
				connMap[requestUser] = conn
				self.connMutex.Unlock()

				sessionNew := Session{
					User:    requestUser,
					Status:  LOGIN,
					AppId:   appid,
					Version: version,
				}
				conn.SetExtraData(sessionNew)
				holmes.Info("user[%s] appid[%d] version[%s] relogin success with conn name[%s]", requestUser, appid, version, conn.GetName())
				self.sendDuobbManagerMsg(requestUser, ext.DUOBB_MANAGER_TYPE_LOGIN)
			}
		} else {
			self.connMutex.Lock()
			//conn.SetName(requestUser + CONN_NAME_DELIMITER + conn.GetName())
			connMap[requestUser] = conn
			self.connMutex.Unlock()

			session := Session{
				User:    requestUser,
				Status:  LOGIN,
				AppId:   appid,
				Version: version,
			}
			conn.SetExtraData(session)
			holmes.Info("user[%s] appid[%d] version[%s] login success with conn name[%s]", requestUser, appid, version, conn.GetName())
			self.sendDuobbManagerMsg(requestUser, ext.DUOBB_MANAGER_TYPE_LOGIN)
		}
	case DUOBB_ACCESS_LOGOUT:
		reqMap, err := self.parseMsg(decodeMsg)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
			break
		}
		user, err := self.checkReqMsgUser(reqMap)
		if err != nil {
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
			break
		}
		if user != requestUser {
			holmes.Error("user[%s] not equal requestUser[%s]", user, requestUser)
			resultResponse.Code = duobb_proto.DUOBB_MSG_LOGOUT_ERROR
			resultResponse.Msg = duobb_proto.MSG_DUOBB_LOGOUT_ERROR
			break
		}
		if c != nil {
			if c.GetName() == conn.GetName() {
				ifLogout = true
				holmes.Info("start to logout user: %s conn: %s", requestUser, c.GetName())
				self.sendDuobbManagerMsg(requestUser, ext.DUOBB_MANAGER_TYPE_LOGOUT)
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
	case DUOBB_ACCESS_HEARTBEAT:
		holmes.Debug("user[%s] appid[%d] conn[%s] in heartbeat", requestUser, appid, conn.GetName())
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
	holmes.Debug("duobb access result msg: %v", resultResponse)
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

func (self *DuobbProcess) sendDuobbManagerMsg(user string, t int) {
	wxMsgReq := ext.DuobbManagerSendMsgReq{
		UserName: user,
	}
	switch t {
	case ext.DUOBB_MANAGER_TYPE_LOGIN:
		wxMsgReq.Msg = fmt.Sprintf(ext.LOGIN_MSG, user)
	case ext.DUOBB_MANAGER_TYPE_LOGOUT:
		wxMsgReq.Msg = fmt.Sprintf(ext.LOGOUT_MSG, user)
	case ext.DUOBB_MANAGER_TYPE_HEALTH_ERROR:
		wxMsgReq.Msg = fmt.Sprintf(ext.HEALTH_ERROR, user)
	}
	self.dmExt.AsyncSendMsg(&ext.DuobbManagerReq{Uri: ext.DUOBB_MANAGER_WECHAT_MSG_URI, Req: wxMsgReq})
}

func (self *DuobbProcess) parseMsg(decodeMsg []byte) (map[string]interface{}, error) {
	request, err := JsonDecode(decodeMsg)
	if err != nil {
		holmes.Error("json decode[%s] error: %v", string(decodeMsg), err)
		return nil, fmt.Errorf("json decode error: %v", err)
	}
	reqMap := request.(map[string]interface{})
	if reqMap == nil {
		holmes.Error("request: %v translate to map error", request)
		return nil, fmt.Errorf("request: %v translate to map error", request)
	}
	return reqMap, nil
}

func (self *DuobbProcess) checkReqMsgUser(reqMap map[string]interface{}) (string, error) {
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

func (self *DuobbProcess) checkReqMsgVersion(reqMap map[string]interface{}) (string, error) {
	version := reqMap["version"]
	if version == nil {
		holmes.Error("reqmap: %v has no version field", reqMap)
		return "", fmt.Errorf("reqmap: %v has no version field", reqMap)
	}
	versionStr, ok := version.(string)
	if !ok {
		holmes.Error("version: %v translate to string error", version)
		return "", fmt.Errorf("version: %v translate to string error", version)
	}
	return versionStr, nil
}

func (self *DuobbProcess) OnConnectCallback(conn tao.Connection) bool {
	holmes.Info("[%s] on connect", conn.GetName())
	return true
}

func (self *DuobbProcess) OnErrorCallback() {
	holmes.Info("on Error")
}

func (self *DuobbProcess) OnCloseCallback(conn tao.Connection) {
	//holmes.Debugln(self.connMap)

	extra := conn.GetExtraData()
	if extra == nil {
		holmes.Error("get extra data error == nil")
		return
	}

	session := extra.(Session)
	if session.Status == LOGOUT_KICKOFF {
		holmes.Info("conn[%s] is kickoff logout clear success.", conn.GetName())
		return
	}

	connName := conn.GetName()
	holmes.Info("conn[%s] session[%v] on close", connName, session)
	//names := strings.Split(conn.GetName(), CONN_NAME_DELIMITER)
	//if len(names) != 2 {
	//	holmes.Error("[%s] name error.", connName)
	//	return
	//}
	self.connMutex.Lock()
	defer self.connMutex.Unlock()
	connMap := self.connMap[session.AppId]
	if connMap == nil {
		holmes.Error("cannot found this appid[%d]", session.AppId)
		return
	}
	c := connMap[session.User]
	if c == nil {
		holmes.Error("[%s] connmap has no this connection.", session.User)
		return
	}
	delete(connMap, session.User)
	holmes.Info("conn[%s] on close clear success.", connName)
}

func (self *DuobbProcess) OnScheduleCallback(now time.Time, data interface{}) {
	conn := data.(tao.Connection)
	extra := conn.GetExtraData()
	var session Session
	if extra == nil {
		holmes.Error("client %d %s get extra data error.", conn.GetNetId(), conn.GetName())
	} else {
		session = extra.(Session)
		if session.Status == LOGOUT_KICKOFF {
			session.CheckLogout++
			if session.CheckLogout == MAX_LOGOUT_KICKOFF {
				holmes.Warn("client %d %s kickoff max times, close it\n", conn.GetNetId(), conn.GetName())
				conn.Close()
				return
			}
			conn.SetExtraData(session)
		}
	}
	lastTimestamp := conn.GetHeartBeat()
	if now.UnixNano()-lastTimestamp > int64(self.cfg.ServerConfig.Timeout)*1e9 {
		holmes.Warn("client %d %s timeout, close it\n", conn.GetNetId(), conn.GetName())
		conn.Close()
		self.sendDuobbManagerMsg(session.User, ext.DUOBB_MANAGER_TYPE_HEALTH_ERROR)
	}
}
