package model

// Tag 标签
type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"unique;not null"`
}

func (*Tag) TableName() (name string) {
	return "tag"
}
