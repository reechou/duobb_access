package ext

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/reechou/duobb_access/config"
	"github.com/reechou/holmes"
)

type DuobbManagerExt struct {
	cfg    *config.Config
	client *http.Client

	stop    chan struct{}
	done    chan struct{}
	msgChan chan *DuobbManagerReq
}

func NewDuobbManagerExt(cfg *config.Config) *DuobbManagerExt {
	dme := &DuobbManagerExt{
		cfg:     cfg,
		client:  &http.Client{},
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
		msgChan: make(chan *DuobbManagerReq, WX_SEND_MSG_CHAN_LEN),
	}
	go dme.run()

	return dme
}

func (self *DuobbManagerExt) Stop() {
	close(self.stop)
	<-self.done
}

func (self *DuobbManagerExt) run() {
	holmes.Debug("duobb manager ext start run.")
	for {
		select {
		case msg := <-self.msgChan:
			self.SendMsg(msg)
		case <-self.stop:
			close(self.done)
			return
		}
	}
}

func (self *DuobbManagerExt) AsyncSendMsg(msg *DuobbManagerReq) {
	holmes.Debug("async duobb manager send msg[%v]", msg)
	select {
	case self.msgChan <- msg:
	case <-time.After(2 * time.Second):
		holmes.Error("async wx send msg timeout.")
		return
	}
}

func (self *DuobbManagerExt) SendMsg(info *DuobbManagerReq) error {
	u := "http://" + self.cfg.DuobbManagerSrv.HostURL + info.Uri
	body, err := json.Marshal(info.Req)
	if err != nil {
		return err
	}
	httpReq, err := http.NewRequest("POST", u, strings.NewReader(string(body)))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	rsp, err := self.client.Do(httpReq)
	defer func() {
		if rsp != nil {
			rsp.Body.Close()
		}
	}()
	if err != nil {
		return err
	}
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	var response DuobbManagerResponse
	err = json.Unmarshal(rspBody, &response)
	if err != nil {
		return err
	}
	if response.Code != DUOBB_MANAGER_RESPONSE_OK {
		holmes.Error("duobb manager send msg error: %v", response)
		return fmt.Errorf("duobb manager send msg error: %s", response.Msg)
	}

	return nil
}
