package main

import (
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hashwing/log"
	"github.com/hashwing/pet-adoption/pkg/common"
	"github.com/hashwing/pet-adoption/pkg/config"
	"github.com/hashwing/pet-adoption/pkg/storage/db"
	"github.com/namsral/flag"
)

var _VERSION_ = "unknown"

// run run in daemonize
func run() {
	confFile := flag.String("c", "/etc/pet-adoption/config.yml", "config file.")
	opt := flag.String("opt", "", "option")
	className := flag.String("class", "", "class name")
	localityName := flag.String("locality", "", "locality name")
	flag.Parse()
	err := config.InitGlobal(*confFile)
	if err != nil {
		log.Error("can't load config file: %v", err)
		panic(err)
	}
	err = db.NewDB(config.Cfg.MysqlConfig.User,
		config.Cfg.MysqlConfig.Password,
		config.Cfg.MysqlConfig.Host,
		config.Cfg.MysqlConfig.DBName)
	if err != nil {
		log.Error("can't init db: %v", err)
		panic(err)
	}

	switch *opt {
	case "add_pet_class":
		if *className == "" {
			log.Error("class name is null")
		}
		p := db.PetClass{
			ID:   common.NewUUID(),
			Name: *className,
		}
		err := db.AddPetClass(p)
		if err != nil {
			log.Error(err)
		}
	case "add_locality":
		provinces, err := db.FindProvinces()
		if err != nil {
			log.Error(err)
			return
		}

		citys, err := db.FindCitys()
		if err != nil {
			log.Error(err)
			return
		}

		pid := common.NewUUID()
		pflag := false
		cid := common.NewUUID()
		cflag := false
		lid := common.NewUUID()

		localityNames := strings.Split(*localityName, "/")
		for _, v := range provinces {
			if v.Name == localityNames[0] {
				pid = v.ID
				pflag = true
				break
			}
		}
		if !pflag {
			db.AddProvince(
				db.Province{
					ID:   pid,
					Name: localityNames[0],
				},
			)
		}

		for _, v := range citys {
			if v.Name == localityNames[1] {
				cid = v.ID
				cflag = true
				break
			}
		}

		if !cflag {
			db.AddCity(
				db.City{
					ID:         cid,
					ProvinceID: pid,
					Name:       localityNames[1],
				},
			)
		}

		db.AddLocality(
			db.Locality{
				ID:     lid,
				CityID: cid,
				Name:   localityNames[2],
			},
		)

	}

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
