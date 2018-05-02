package dialect

import "database/sql"

// MySQL is the dialect used for MySQL and compatible databases, such as MariaDB.
var MySQL mysql

type mysql struct{}

func (mysql) TxOptions() *sql.TxOptions {
	return &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	}
}
