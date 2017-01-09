package duobb

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/reechou/duobb_access/models"
	"github.com/reechou/holmes"
)

type DuobbAccessCfgHttpRes struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func (self *DuobbProcess) httpInit() {
	http.HandleFunc("/cfg/create_service", self.CreateService)
	http.HandleFunc("/cfg/create_service_method", self.CreateServiceMethod)
	http.HandleFunc("/cfg/load_service", self.LoadService)

	http.HandleFunc("/duobb/login_user_num", self.GetDuobbLoginUserNum)
	http.HandleFunc("/duobb/create_push_msg", self.CreatePushMsg)

	holmes.Info("duobb access cfg http listen on[%s]", self.cfg.CfgHost)
	if err := http.ListenAndServe(self.cfg.CfgHost, nil); err != nil {
		holmes.Errorln(err)
		return
	}
}

func (self *DuobbProcess) CreateService(w http.ResponseWriter, r *http.Request) {
	req := &models.Service{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateService json decode error: %v", err)
		return
	}

	rsp := &DuobbAccessCfgHttpRes{Code: 0}
	err := models.CreateService(req)
	if err != nil {
		holmes.Error("CreateService error: %v", err)
		rsp.Code = 1
		rsp.Msg = err.Error()
	}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *DuobbProcess) CreateServiceMethod(w http.ResponseWriter, r *http.Request) {
	req := &models.ServiceMethod{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateServiceMethod json decode error: %v", err)
		return
	}

	rsp := &DuobbAccessCfgHttpRes{Code: 0}
	err := models.CreateServiceMethod(req)
	if err != nil {
		holmes.Error("CreateServiceMethod error: %v", err)
		rsp.Code = 1
		rsp.Msg = err.Error()
	}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *DuobbProcess) LoadService(w http.ResponseWriter, r *http.Request) {
	rsp := &DuobbAccessCfgHttpRes{Code: 0}
	err := self.initService()
	if err != nil {
		holmes.Error("LoadService error: %v", err)
		rsp.Code = 1
		rsp.Msg = err.Error()
	}
	WriteJSON(w, http.StatusOK, rsp)
}

func (self *DuobbProcess) GetDuobbLoginUserNum(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	appid := 0
	if len(r.Form["appid"]) != 0 {
		appidStr := r.Form["appid"][0]
		appid, _ = strconv.Atoi(appidStr)
	}

	type DuobbLoginUser struct {
		User string `json:"user"`
		IP   string `json:"ip"`
	}
	type DuobbLoginUsers struct {
		Count int              `json:"count"`
		Users []DuobbLoginUser `json:"users"`
	}

	rsp := &DuobbAccessCfgHttpRes{Code: 0}
	//rsp.Data = self.server.GetConnectionMap().Size()
	var users DuobbLoginUsers
	self.connMutex.Lock()
	connMap := self.connMap[appid]
	for k, v := range connMap {
		user := DuobbLoginUser{
			User: k,
			IP:   v.GetName(),
		}
		users.Users = append(users.Users, user)
	}
	self.connMutex.Unlock()
	users.Count = len(users.Users)
	rsp.Data = users

	WriteJSON(w, http.StatusOK, rsp)
}

func (self *DuobbProcess) CreatePushMsg(w http.ResponseWriter, r *http.Request) {
	req := &models.DuobbPushMsg{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		holmes.Error("CreateServiceMethod json decode error: %v", err)
		return
	}

	rsp := &DuobbAccessCfgHttpRes{Code: 0}
	req.Msg = url.QueryEscape(req.Msg)
	holmes.Debug("create push msg req: %v", req)
	err := models.CreateDuobbPushMsg(req)
	if err != nil {
		holmes.Error("CreatePushMsg error: %v", err)
		rsp.Code = 1
		rsp.Msg = err.Error()
	}

	WriteJSON(w, http.StatusOK, rsp)
}

func WriteJSON(w http.ResponseWriter, code int, v interface{}) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "x-requested-with,content-type")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}
