package main

import (
	"fmt"
	"os"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hashwing/log"
	"github.com/hashwing/pet-adoption/pkg/config"
	"github.com/hashwing/pet-adoption/pkg/storage/db"
	_ "github.com/hashwing/pet-adoption/routers"
	"github.com/namsral/flag"
)

var _VERSION_ = "unknown"

// run run in daemonize
func run() {
	confFile := flag.String("c", "/etc/pet-adoption/config.yml", "config file.")
	flag.Parse()
	log.Info("init config...")
	err := config.InitGlobal(*confFile)
	if err != nil {
		log.Error("can't load config file: %v", err)
		panic(err)
	}

	log.Info("init db...")
	err = db.NewDB(config.Cfg.MysqlConfig.User,
		config.Cfg.MysqlConfig.Password,
		config.Cfg.MysqlConfig.Host,
		config.Cfg.MysqlConfig.DBName)
	if err != nil {
		log.Error("can't init db: %v", err)
		panic(err)
	}

	beego.BConfig.AppName = "pet-adoption"
	beego.BConfig.RunMode = beego.PROD
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.Log.FileLineNum = true
	beego.BConfig.Listen.HTTPPort = config.Cfg.Port
	log.Info("listen in %d", config.Cfg.Port)
	beego.Run()
}

func main() {
	version()
	run()
}

func version() {
	args := os.Args
	if len(args) > 1 {
		if args[1] == "-v" || args[1] == "--version" {
			fmt.Println(_VERSION_)
			os.Exit(0)
		}
	}
}
