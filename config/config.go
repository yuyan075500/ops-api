package config

import "time"

const (
	ListenAddr   = "0.0.0.0:8000"
	DBHost       = "127.0.0.1"
	DBPort       = 3306
	DBName       = "ops"
	DBUser       = "root"
	DBPassword   = "p-0p-0p-0"
	MaxIdleConns = 10
	MaxOpenConns = 100
	MaxLifeTime  = 30 * time.Second
)
