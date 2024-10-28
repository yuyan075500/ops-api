package service

import (
	"ops-api/dao"
	"ops-api/global"
	"ops-api/middleware"
	"ops-api/model"
)

var Task task

type task struct{}

// TaskCreate 创建构体
type TaskCreate struct {
	Name          string `json:"name" binding:"required"`
	Type          uint   `json:"type" binding:"required"`
	CronExpr      string `json:"cron_expr" binding:"required"`
	BuiltInMethod string `json:"built_in_method" binding:"required"`
	Enabled       *bool  `json:"enabled" binding:"required"`
}

// AddTask 创建定时任务
func (t *task) AddTask(data *TaskCreate) (err error) {

	task := &model.ScheduledTask{
		Name:          data.Name,
		Type:          data.Type,
		CronExpr:      data.CronExpr,
		BuiltInMethod: data.BuiltInMethod,
		Enabled:       *data.Enabled,
	}

	// 数据库中新增任务
	if err := dao.Task.AddTask(task); err != nil {
		return err
	}

	// 如果任务未启用，则停止任务，否则更新任务
	if task.Enabled {
		if err := middleware.AddOrUpdateTask(*task); err != nil {
			return err
		}
	}

	return nil
}

// DeleteTask 删除定时任务
func (t *task) DeleteTask(id int) (err error) {

	// 查询删除的任务
	task := &model.ScheduledTask{}
	if err := global.MySQLClient.Where("id = ?", id).First(task).Error; err != nil {
		return err
	}

	// 删除已经加载的任务
	if task.EntryID != nil {
		global.CornSchedule.Remove(*task.EntryID)
	}

	// 删除任务本身
	if err = dao.Task.DeleteTask(id); err != nil {
		return err
	}
	return nil
}

// UpdateTask 更新定时任务
func (t *task) UpdateTask(data *dao.TaskUpdate) error {

	// 更新任务本身
	if err := dao.Task.UpdateTask(data); err != nil {
		return err
	}

	// 查询更新的任务
	task := &model.ScheduledTask{}
	if err := global.MySQLClient.Where("id = ?", data.ID).First(task).Error; err != nil {
		return err
	}

	// 如果任务未启用，则停止任务，否则更新任务
	if task.Enabled {
		if err := middleware.AddOrUpdateTask(*task); err != nil {
			return err
		}
	} else {
		if task.EntryID != nil {
			global.CornSchedule.Remove(*task.EntryID)
		}
	}

	return nil
}

// GetTaskList 获取定时任务列表
func (t *task) GetTaskList(name string, page, limit int) (data *dao.TaskList, err error) {
	data, err = dao.Task.GetTaskList(name, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}
