package entity

import "time"

type Token struct {
	Tid  uint64 `json:"tid" gorm:"primarykey"`
	Uid  uint64 `json:"uid"`
	Uuid string
	Exp  time.Time
}
