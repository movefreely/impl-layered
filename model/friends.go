package model

type Friend struct {
	Id      int64  `json:"id" gorm:"primary_key"`
	OneId   uint64 `json:"one_id"`
	OtherId uint64 `json:"other_id"`
	Status  int    `json:"status" gorm:"default:0; type:tinyint"` // 0: pending, 1: accepted, 2: rejected
}
