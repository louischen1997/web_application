package mdl

type Dbtable struct {
	First_name      string `gorm:"type:varchar(45)"`
	Last_name       string `gorm:"type:varchar(45)"`
	Password        string `gorm:"type:varchar(225)"`
	Username        string `gorm:"type:varchar(45)"`
	ID              string `gorm:"type:varchar(45);primary_key"`
	Account_created string `gorm:"type:varchar(100)"`
	Account_updated string `gorm:"type:varchar(100)"`
}

type Dbtables []Dbtable
