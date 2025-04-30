package entity

type Vds struct {
	Vid    uint64 `json:"uid" gorm:"primarykey"`
	Uid    uint64 `json:"Uid"` // список пользователей
	Name   string //название машины
	Img    string
	Ip     string
	Login  string
	Pass   string
	Status bool //вкл/выкл
}
