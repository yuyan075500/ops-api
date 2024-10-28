package dao

import (
	"ops-api/global"
	"ops-api/model"
)

var Task task

type task struct{}

// TaskList 返回给前端列表结构体
type TaskList struct {
	Items []*model.ScheduledTask `json:"items"`
	Total int64                  `json:"total"`
}

// TaskUpdate 更新构体
type TaskUpdate struct {
	ID            uint   `json:"id" binding:"required"`
	Name          string `json:"name"`
	Type          uint   `json:"type"`
	CronExpr      string `json:"cron_expr"`
	Method        uint   `json:"method"`
	BuiltInMethod string `json:"built_in_method"`
	Enabled       bool   `json:"enabled"`
}

// AddTask 新增定时任务
func (t *task) AddTask(data *model.ScheduledTask) (err error) {
	if err := global.MySQLClient.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

// DeleteTask 删除定时任务
func (t *task) DeleteTask(id int) (err error) {
	if err := global.MySQLClient.Where("id = ?", id).Unscoped().Delete(&model.ScheduledTask{}).Error; err != nil {
		return err
	}
	return nil
}

// UpdateTask 修改定时任务
func (t *task) UpdateTask(data *TaskUpdate) (err error) {
	if err := global.MySQLClient.Model(&model.ScheduledTask{}).Select("*").Where("id = ?", data.ID).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// GetTaskList 获取定时任务列表
func (t *task) GetTaskList(name string, page, limit int) (data *TaskList, err error) {

	// 定义数据的起始位置
	startSet := (page - 1) * limit

	// 定义返回的内容
	var (
		items []*model.ScheduledTask
		total int64
	)

	// 获取菜单列表
	tx := global.MySQLClient.Model(&model.ScheduledTask{}).
		Where("name like ?", "%"+name+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Find(&items)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &TaskList{
		Items: items,
		Total: total,
	}, nil
}
