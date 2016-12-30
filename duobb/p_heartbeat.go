package duobb

import (
	"fmt"
	"net/url"

	"github.com/reechou/duobb_access/models"
	"github.com/reechou/holmes"
)

func (self *DuobbProcess) heartbeatKickoff(user string) (interface{}, error) {
	type PushMsgKickoff struct {
		User    string `json:"user"`
		Because string `json:"because"`
	}
	msg := &PushMsgKickoff{
		User:    user,
		Because: "该账号已在别处登录",
	}
	msgBytes, err := JsonEncode(msg)
	if err != nil {
		return nil, err
	}
	pushMsg := models.DuobbPushMsg{
		Type: 3,
		Msg:  url.QueryEscape(string(msgBytes)),
	}
	return []models.DuobbPushMsg{pushMsg}, nil
}

func (self *DuobbProcess) checkHeartbeat(decodeMsg []byte) (interface{}, error) {
	request, err := JsonDecode(decodeMsg)
	if err != nil {
		holmes.Error("json decode[%s] error: %v", string(decodeMsg), err)
		return nil, err
	}
	reqMap := request.(map[string]interface{})
	if reqMap == nil {
		holmes.Error("request: %v translate to map error", request)
		return nil, fmt.Errorf("request: %v translate to map error", request)
	}
	lastPushMsgTime := reqMap["lastPushMsgTime"]
	if lastPushMsgTime == nil {
		holmes.Error("reqmap: %v has no lastPushMsgTime field", reqMap)
		return nil, fmt.Errorf("reqmap: %v has no lastPushMsgTime field", reqMap)
	}
	t, ok := lastPushMsgTime.(float64)
	if !ok {
		holmes.Error("lastPushMsgTime: %v translate to float64 error", lastPushMsgTime)
		return nil, fmt.Errorf("lastPushMsgTime: %v translate to float64 error", lastPushMsgTime)
	}
	pushMsgList, err := models.GetDuobbPushMsgFromTime(int64(t))
	if err != nil {
		holmes.Error("get push msg error: %v", err)
		return nil, fmt.Errorf("get push msg error: %v", err)
	}
	if len(pushMsgList) == 0 {
		return nil, nil
	}

	return pushMsgList, nil
}
