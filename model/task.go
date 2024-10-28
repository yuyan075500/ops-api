package model

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"time"
)

type ScheduledTask struct {
	gorm.Model
	Name           string        `json:"name"`                                // 任务名称
	Type           uint          `json:"type"`                                // 任务类型（1：间隔任务，2：周期任务）
	CronExpr       string        `json:"cron_expr" gorm:"default:null"`       // 周期任务的表达式
	Method         uint          `json:"method"`                              // 执行方法（1：调用URL，2：内置方法）
	BuiltInMethod  string        `json:"built_in_method" gorm:"default:null"` // 内置方法
	Enabled        bool          `json:"enabled"`                             // 是否启用
	LastRunAt      *time.Time    `json:"last_run_at"`                         // 上次运行时间
	NextRunAt      *time.Time    `json:"next_run_at"`                         // 下次运行时间
	EntryID        *cron.EntryID `json:"entry_id"`                            // 任务运行ID
	ExecutionCount uint          `json:"execution_count" gorm:"default:0"`    // 执行次数
}
