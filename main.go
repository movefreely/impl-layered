package main

import (
	"in-server/controller"
	"in-server/util/sqlinit"
)

func main() {
	sqlinit.SqlInit()
	controller.Entrance()
	//err := sqlinit.DB.AutoMigrate(&model.Friend{})
	//err = sqlinit.DB.Table("latest_message").AutoMigrate(&model.Message{})
	//if err != nil {
	//	return
	//}
	//if err != nil {
	//	fmt.Println("创建数据表失败： ", err)
	//}
}
