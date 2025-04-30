package entity

type User struct {
	Uid    uint64 `json:"uid" gorm:"primarykey"`
	Login  string
	Pass   string
	Name   string
	Access bool
}
