package controller

import (
	"os"
	"runtime"

	"github.com/reechou/duobb_access/config"
	"github.com/reechou/duobb_access/models"
)

type Logic struct {
	cfg    *config.Config
	server *DuobbAccessServer
}

func NewLogic(cfg *config.Config) *Logic {
	models.InitDB(cfg)

	l := &Logic{
		cfg:    cfg,
		server: NewDuobbAccessServer(cfg),
	}
	if cfg.Debug {
		EnableDebug()
	}

	return l
}

func (self *Logic) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	self.server.Start()
}

func EnableDebug() {
	os.Setenv("DEBUG", "1")
}
