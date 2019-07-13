package db

import (
	"fmt"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// MysqlDB mysql engine
var MysqlDB *xorm.Engine

// NewDB new db
func NewDB(user, passwd, url, dbName string) error {
	var err error
	dbURL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		user,
		passwd,
		url,
		dbName)
	MysqlDB, err = xorm.NewEngine("mysql", dbURL)
	if err != nil {
		return err
	}

	err = MysqlDB.Sync2(new(Province), new(City), new(Locality), new(PetClass), new(PetPublic), new(AdoptionApply), new(User))
	return err
}
