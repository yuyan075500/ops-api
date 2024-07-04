package utils

import (
	"errors"
	"github.com/go-sql-driver/mysql"
)

// IsForeignKeyConstraintError 判断是否为外建关联无法删除错误
func IsForeignKeyConstraintError(err error) bool {

	var mysqlError *mysql.MySQLError

	if errors.As(err, &mysqlError) {
		return mysqlError.Number == 1451
	}

	return false
}

// IsDuplicateEntryError 关断是否为Duplicate entry错误
func IsDuplicateEntryError(err error) bool {

	var mysqlError *mysql.MySQLError

	if errors.As(err, &mysqlError) {
		return mysqlError.Number == 1062
	}

	return false
}
