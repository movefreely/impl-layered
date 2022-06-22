package model

import (
	"time"
)

type User struct {
	ID       uint64    `gorm:"primary_key" json:"id"`                     // 用户ID
	Email    string    `gorm:"type:varchar(60);not null" json:"email"`    // 邮箱
	Password string    `gorm:"type:varchar(60);not null" json:"password"` // 密码
	Avatar   string    `gorm:"type:varchar(100)" json:"avatar"`           // 头像
	Sex      string    `gorm:"type:varchar(5)" json:"sex"`                // 性别
	Nickname string    `gorm:"type:varchar(30)" json:"nickname"`          // 昵称
	Online   string    `gorm:"type:varchar(10)" json:"online"`            // 是否在线 1为在线 0为不在线
	CreateAt time.Time `json:"create_at"`                                 // 创建时间
}
