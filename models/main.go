package models



import (

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var  db  *gorm.DB

func init() {

	db, _ = gorm.Open("mysql", "root:lost52wf@/testforfox?charset=utf8&parseTime=True&loc=Local")
	//defer db.Close()
	//return db

}