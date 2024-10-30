package service

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/wonderivan/logger"
	"ops-api/dao"
	"ops-api/global"
	"ops-api/model"
	"time"
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
		if err := AddOrUpdateTask(*task); err != nil {
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
		if err := AddOrUpdateTask(*task); err != nil {
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

// GetTaskLogList 获取定时任务执行日志列表
func (t *task) GetTaskLogList(id uint, page, limit int) (data *dao.TaskLogList, err error) {
	data, err = dao.Task.GetTaskLogList(id, page, limit)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TaskInit 初始化任务调度器并加载数据库中的任务
func TaskInit() error {
	global.CornSchedule = cron.New(cron.WithChain())

	// 加载数据库中的启用任务
	var tasks []model.ScheduledTask
	if err := global.MySQLClient.Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return errors.New(fmt.Sprintf("failed to load tasks: %v", err))
	}

	// 遍历任务并将其添加到调度器
	for _, task := range tasks {
		if err := AddOrUpdateTask(task); err != nil {
			return errors.New(fmt.Sprintf("failed to add task %s to cron, %v", task.Name, err))
		}
	}

	// 启动调度器
	global.CornSchedule.Start()

	// 打印日志
	logger.Info("定时任务初始化成功.")

	return nil
}

// AddOrUpdateTask 添加或更新单个任务到调度器
func AddOrUpdateTask(task model.ScheduledTask) error {

	var entryID cron.EntryID

	// 检查任务是否已经运行
	if task.EntryID != nil {
		global.CornSchedule.Remove(*task.EntryID)
	}

	// 使用任务的 Cron 表达式
	entryID, err := global.CornSchedule.AddFunc(task.CronExpr, func() {

		// 创建执行日志记录
		startTime := time.Now()
		execLog := model.ScheduledTaskExecLog{
			ScheduledTaskID: task.ID,
			RunAt:           &startTime,
		}
		if err := global.MySQLClient.Create(&execLog).Error; err != nil {
			return
		}

		// 执行任务逻辑
		if task.Type == 1 {
			executeURLTask()
		} else {
			executeBuiltInMethod(task, &execLog)
		}

		// 更新任务信息
		global.MySQLClient.Model(&task).Updates(map[string]interface{}{
			"last_run_at":     startTime,
			"next_run_at":     global.CornSchedule.Entry(entryID).Next,
			"entry_id":        entryID,
			"execution_count": task.ExecutionCount + 1,
		})
	})

	if err != nil {
		return err
	}

	return nil
}

// executeURLTask 请求URL
func executeURLTask() {}

// executeBuiltInMethod 内置方法调用
func executeBuiltInMethod(task model.ScheduledTask, execLog *model.ScheduledTaskExecLog) {

	defer func() {
		finishTime := time.Now()
		global.MySQLClient.Model(execLog).Updates(map[string]interface{}{
			"finish_at": finishTime,
		})
	}()

	// 用户同步
	if task.BuiltInMethod == "user_sync" {
		if err := AD.LDAPUserSync(); err != nil {
			global.MySQLClient.Model(execLog).Update("result", err.Error())
			global.MySQLClient.Model(&task).Update("LastRunResult", "失败")
			logger.Warn("任务执行失败:", err.Error())
		} else {
			global.MySQLClient.Model(execLog).Update("result", "成功")
			global.MySQLClient.Model(&task).Update("LastRunResult", "成功")
		}

	}
}
