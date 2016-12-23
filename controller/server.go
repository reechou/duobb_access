package controller

import (
	"time"

	"github.com/reechou/duobb_access/config"
	"github.com/reechou/duobb_access/duobb"
	"github.com/reechou/holmes"
	"github.com/reechou/tao"
)

type DuobbAccessServer struct {
	tao.Server
	process *duobb.DuobbProcess
	cfg     *config.Config
}

func NewDuobbAccessServer(cfg *config.Config) *DuobbAccessServer {
	das := &DuobbAccessServer{
		Server:  tao.NewTCPServer(cfg.Host),
		process: duobb.NewDuobbProcess(cfg),
		cfg:     cfg,
	}
	das.init()

	return das
}

func (self *DuobbAccessServer) init() {
	//tao.MonitorOn(12345)
	tao.Register(duobb.DuobbMsgCMD, self.process.DeserializeMessage, self.process.ProcessDuobbMessage)

	self.Server.SetOnConnectCallback(self.process.OnConnectCallback)
	self.Server.SetOnErrorCallback(self.process.OnErrorCallback)
	self.Server.SetOnCloseCallback(self.process.OnCloseCallback)
	self.Server.SetOnScheduleCallback(time.Duration(self.cfg.ServerConfig.CheckTimeoutInterval)*time.Second, self.process.OnScheduleCallback)
}

func (self *DuobbAccessServer) Start() {
	holmes.Info("duobb access listen on[%s]", self.cfg.Host)
	self.Server.Start()
}

func (self *DuobbAccessServer) Stop() {
	self.Server.Close()
}
