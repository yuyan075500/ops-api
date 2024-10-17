package service

import (
	"ops-api/dao"
	"ops-api/model"
)

var Tag tag

type tag struct{}

// GetTagList 获取标签列表
func (t *tag) GetTagList(name string) (data []*model.Tag, err error) {
	data, err = dao.Tag.GetTagList(name)
	if err != nil {
		return nil, err
	}
	return data, nil
}
