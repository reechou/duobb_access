package duobb

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strings"
	"sync"

	"github.com/reechou/duobb_access/config"
	"github.com/reechou/duobb_access/models"
	"github.com/reechou/duobb_access/ext"
	"github.com/reechou/duobb_proto"
	"github.com/reechou/holmes"
	"github.com/reechou/tao"
)

type DuobbProcess struct {
	smMutex   sync.Mutex
	connMutex sync.Mutex

	serviceMap       map[string]*BackendService
	serviceMethodMap map[string]*BackendServiceMethod
	connMap          map[int]map[string]tao.Connection

	cfg        *config.Config
	sec        *Security
	jsonRpc    *JsonRpc
	normalHttp *NormalHttp
	dmExt      *ext.DuobbManagerExt
	server     tao.Server
}

func NewDuobbProcess(cfg *config.Config, server tao.Server) *DuobbProcess {
	dp := &DuobbProcess{
		cfg:              cfg,
		sec:              &Security{},
		jsonRpc:          NewJsonRpc(),
		normalHttp:       NewNormalHttp(),
		dmExt:            ext.NewDuobbManagerExt(cfg),
		serviceMap:       make(map[string]*BackendService),
		serviceMethodMap: make(map[string]*BackendServiceMethod),
		connMap:          make(map[int]map[string]tao.Connection),
		server:           server,
	}
	for _, v := range apps {
		dp.connMap[v] = make(map[string]tao.Connection)
	}
	
	dp.initService()
	go dp.httpInit()

	return dp
}

func (self *DuobbProcess) initService() error {
	serviceList, err := models.LoadService()
	if err != nil {
		holmes.Error("load service error: %v", err)
		return err
	}
	for _, v := range serviceList {
		hosts := strings.Split(v.Hosts, ",")
		self.RegisterService(v.ServiceName, hosts, v.Type)
	}
	serviceMethodList, err := models.LoadServiceMethod()
	if err != nil {
		holmes.Error("load service method error: %v", err)
		return err
	}
	for _, v := range serviceMethodList {
		self.RegisterServiceMethod(v.ServiceMethod, v.Uri)
	}
	return nil
	//bs := strings.Split(self.cfg.Backend.BackendHosts, ",")
	//for _, v := range bs {
	//	vs := strings.Split(v, "*")
	//	if len(vs) != 3 {
	//		holmes.Error("backend host[%s] error: %v", v)
	//		continue
	//	}
	//	hosts := strings.Split(vs[1], "&")
	//	self.RegisterService(vs[0], hosts, vs[2])
	//}
	//self.RegisterService("DuobbAccountService", []string{"127.0.0.1:7878"})
}

func (self *DuobbProcess) DeserializeMessage(data []byte) (tao.Message, error) {
	if data == nil {
		return nil, ErrorNilData
	}
	dataLen := uint32(len(data))

	buffer := bytes.NewBuffer(data)
	
	var appid int32
	binary.Read(buffer, binary.LittleEndian, &appid)
	dataLen -= 4
	
	var len uint32
	binary.Read(buffer, binary.LittleEndian, &len)
	dataLen -= 4
	if len > dataLen {
		return nil, ErrorIllegalData
	}
	userNameBytes := make([]byte, len)
	binary.Read(buffer, binary.LittleEndian, userNameBytes)
	dataLen -= len

	binary.Read(buffer, binary.LittleEndian, &len)
	dataLen -= 4
	if len > dataLen {
		return nil, ErrorIllegalData
	}
	methodBytes := make([]byte, len)
	binary.Read(buffer, binary.LittleEndian, methodBytes)
	dataLen -= len

	binary.Read(buffer, binary.LittleEndian, &len)
	dataLen -= 4
	if len > dataLen {
		return nil, ErrorIllegalData
	}
	msgBytes := make([]byte, len)
	binary.Read(buffer, binary.LittleEndian, msgBytes)

	msg := &DuobbMsg{
		AppId:    appid,
		UserName: userNameBytes,
		Method:   methodBytes,
		Msg:      msgBytes,
	}
	return msg, nil
}

func (self *DuobbProcess) ProcessDuobbMessage(ctx tao.Context, conn tao.Connection) {
	msg := ctx.Message().(*DuobbMsg)
	if msg == nil {
		holmes.Error("error duobb msg: %v.", msg)
		return
	}
	holmes.Debug("user[%s] method[%s] msglen[%d]", string(msg.UserName), string(msg.Method), len(msg.Msg))

	var secretKey []byte

	rsp := &DuobbMsg{
		AppId:    msg.AppId,
		UserName: msg.UserName,
		Method:   msg.Method,
	}
	var resultMsg []byte
	account, err := self.GetDuobbAccount(string(msg.UserName))
	if err != nil {
		holmes.Error("get duobb account[%s] error: %v", string(msg.UserName), err)
		msgRsp := &duobb_proto.Response{
			Code: duobb_proto.DUOBB_MSG_GET_ACCOUNT_ERROR,
			Msg:  err.Error(),
		}
		result, err := JsonEncode(msgRsp)
		if err != nil {
			holmes.Error("json encode[%v] error: %v", msgRsp, err)
			return
		}
		resultMsg = result
	} else {
		secretKey = self.getSecretKey(account)
		decodeMsg, err := self.decodeMsg(secretKey, msg.Msg)
		if err != nil {
			holmes.Error("decode msg error: %v", err)
			return
		} else {
			if !self.checkMsg(msg, secretKey, decodeMsg, conn) {
				return
			}
			processResult, err := self.process(string(msg.Method), decodeMsg)
			if err != nil {
				holmes.Error("backend process msg[%s] error: %v", string(decodeMsg), err)
				msgRsp := &duobb_proto.Response{
					Code: duobb_proto.DUOBB_MSG_PROCESS_ERROR,
					Msg:  err.Error(),
				}
				result, err := JsonEncode(msgRsp)
				if err != nil {
					holmes.Error("json encode[%v] error: %v", msgRsp, err)
					return
				}
				resultMsg = result
			} else {
				resultMsg = processResult
			}
		}
	}
	resultMsg, err = self.encodeMsg(secretKey, resultMsg)
	if err != nil {
		holmes.Error("user[%s] response msg encode error: %v", string(msg.UserName), err)
		return
	}
	rsp.Msg = resultMsg
	err = conn.Write(rsp)
	if err != nil {
		holmes.Error("conn[%s] write response msg[%v] error: %v", conn.GetName(), rsp, err)
	}
	holmes.Debug("duobb msg from backend back to front success.")
}

func (self *DuobbProcess) encodeMsg(secretKey, msg []byte) ([]byte, error) {
	msgEncode1 := self.sec.Base64Encode(msg)
	msgEncode1 = append(secretKey, msgEncode1...)
	//return self.sec.Base64Encode(msgEncode1)
	return self.sec.GzipEncode(msgEncode1)
}

func (self *DuobbProcess) decodeMsg(secretKey, msg []byte) ([]byte, error) {
	//msgDecode1, err := self.sec.Base64Decode(string(msg))
	msgDecode1, err := self.sec.GzipDecode(msg)
	if err != nil {
		return nil, err
	}
	msgDecode2 := bytes.Replace(msgDecode1, secretKey, []byte(""), -1)

	return self.sec.Base64Decode(string(msgDecode2))
}

func (self *DuobbProcess) getSecretKey(account *models.DuobbAccount) []byte {
	return self.sec.Md5Of32(self.sec.Md5Of32([]byte(account.Password + account.UserName)))
}

func JsonDecode(v []byte) (interface{}, error) {
	var f interface{}
	err := json.Unmarshal(v, &f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func JsonEncode(v interface{}) ([]byte, error) {
	r, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return r, nil
}
