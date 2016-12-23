package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

type AccountDBInfo struct {
	User   string
	Pass   string
	Host   string
	DBName string
}

type ServerConfig struct {
	CheckTimeoutInterval int
	Timeout              int
}

type Config struct {
	Debug   bool
	Path    string
	Logging bool
	PidFile string
	Version string

	Host    string
	CfgHost string

	AccountDBInfo
	ServerConfig
}

func NewConfig() *Config {
	c := new(Config)
	initFlag(c)

	if c.Path == "" {
		fmt.Println("duobb_access must run with config file, please check.")
		os.Exit(0)
	}

	cfg, err := ini.Load(c.Path)
	if err != nil {
		fmt.Printf("ini[%s] load error: %v\n", c.Path, err)
		os.Exit(1)
	}
	cfg.BlockMode = false
	err = cfg.MapTo(c)
	if err != nil {
		fmt.Printf("config MapTo error: %v\n", err)
		os.Exit(1)
	}

	return c
}

func initFlag(c *Config) {
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	v := fs.Bool("v", false, "Print version and exit")
	fs.StringVar(&c.Path, "c", "", "duobb access server config file.")

	fs.Parse(os.Args[1:])
	fs.Usage = func() {
		fmt.Println("Usage: " + os.Args[0] + " -c api.ini")
		fmt.Printf("\nglobal flags:\n")
		fs.PrintDefaults()
	}

	if *v {
		fmt.Println("version: 0.0.1")
		os.Exit(0)
	}
}
