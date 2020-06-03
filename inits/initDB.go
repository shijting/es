package inits

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)
var db *gorm.DB

func init() {
	var err error
	//a := "			  mysql57:shijinting0510@tcp(106.53.5.146:3306)/edu?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql",
		"mysql57:shijinting0510@tcp(106.53.5.146:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(10)

}
func  GetDB() *gorm.DB {
	return db
}
