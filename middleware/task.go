package middleware

import (
	"errors"
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/wonderivan/logger"
	"ops-api/global"
	"ops-api/model"
	"time"
)

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

		// 执行任务逻辑
		if task.Method == 1 {
			executeURLTask()
		} else {
			executeBuiltInMethod(task.BuiltInMethod)
		}

		// 更新任务信息
		global.MySQLClient.Model(&task).Updates(map[string]interface{}{
			"last_run_at":     time.Now(),
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
func executeBuiltInMethod(method string) {
	logger.Info("执行内置方法: " + method)
}
