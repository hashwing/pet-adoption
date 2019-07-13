package db_test

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/hashwing/pet-adoption/pkg/storage/db"
)

func TestFindAdoptionPublics(t *testing.T) {
	err := db.NewDB("root", "", "127.0.0.1:3306", "pet")
	if err != nil {
		t.Error(err)
		return
	}
	//db.FindAdoptionPublics()
}
