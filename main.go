package main

import (
	"github.com/reechou/duobb_access/config"
	"github.com/reechou/duobb_access/controller"
	"github.com/reechou/holmes"
)

func main() {
	defer holmes.Start(holmes.LogFilePath("./log"),
		holmes.EveryDay,
		holmes.AlsoStdout,
		holmes.DebugLevel).Stop()
	controller.NewLogic(config.NewConfig()).Run()
}
