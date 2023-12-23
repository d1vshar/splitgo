package model

import config "github.com/d1vshar/splitgo/config"

type Category struct {
	Meta config.BaseModel `gorm:"embedded" json:"meta"`
	Name string           `json:"name"`
}

func (m *Category) TableName() string {
	return "public._categories"
}
