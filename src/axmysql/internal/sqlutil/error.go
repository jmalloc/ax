package sqlutil

import (
	"github.com/go-sql-driver/mysql"
)

const (
	mysqlDupEntry = 1062 // https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html#error_er_dup_entry
	mysqlDeadLock = 1213 // https://dev.mysql.com/doc/refman/5.5/en/error-messages-server.html#error_er_lock_deadlock
)

// IsDuplicateEntry returns true if err represents a MySQL duplicate entry error.
func IsDuplicateEntry(err error) bool {
	e, ok := err.(*mysql.MySQLError)
	return ok && e.Number == mysqlDupEntry
}

// IsDeadlock returns true if err represents a MySQL deadlock condition.
func IsDeadlock(err error) bool {
	e, ok := err.(*mysql.MySQLError)
	return ok && e.Number == mysqlDeadLock
}
