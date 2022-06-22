package model

import "time"

type Message struct {
	Id       uint64    `json:"id" form:"id" gorm:"primary key"`                 //消息ID
	FromId   uint64    `json:"from_id" form:"from"`                             //谁发的
	ToId     uint64    `json:"to_id" form:"to"`                                 //对端用户ID/群ID
	Type     uint8     `json:"type" form:"cmd"`                                 //群聊还是私聊：定义0为私聊, 1为群聊, 2为系统消息
	Media    int       `json:"media" form:"media" gorm:"type:int"`              //消息按照什么样式展示 0为文字，1为图片，2为语音，3为视频，4为文件
	Content  string    `json:"content" form:"content" gorm:"type:varchar(300)"` //消息的内容(文字)
	Url      string    `json:"url" form:"url" gorm:"type:varchar(300)"`         //url(图片，语音，视频，文件)
	IsAd     bool      `json:"is_ad" form:"is_ad"`                              //是否是广告
	CreateAt time.Time `json:"create_at" form:"create_at"`                      //消息发送时间
}
