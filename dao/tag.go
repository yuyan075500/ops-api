package dao

import (
	"gorm.io/gorm"
	"ops-api/global"
	"ops-api/model"
)

var Tag tag

type tag struct{}

// GetTagList 获取标签列表
func (t *tag) GetTagList(name string) (data []*model.Tag, err error) {

	var tags []*model.Tag

	tx := global.MySQLClient.Model(&model.Tag{}).
		Find(&tags)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return tags, nil
}

// FirstCreateTag 创建标签
func (t *tag) FirstCreateTag(tx *gorm.DB, name string) (data *model.Tag, err error) {

	var tag model.Tag

	// 如果不存在则创建
	if err := tx.Where("name = ?", name).Attrs(model.Tag{Name: name}).FirstOrCreate(&tag).Error; err != nil {
		return nil, err
	}

	return &tag, nil
}
