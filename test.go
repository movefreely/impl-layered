package main

import (
	"in-server/model"
	"in-server/util/sqlinit"
	"time"
)

func main1() {
	sqlinit.SqlInit()
	//msg := service.Message{}
	message := model.Message{
		Type:     0,
		FromId:   12234,
		ToId:     12324,
		Content:  "vfdjn34",
		CreateAt: time.Now(),
	}
	//err := msg.UpdateMessage(&message)
	//if err != nil {
	//	fmt.Println(err)
	//}
	sqlinit.DB.Create(&message)
}
