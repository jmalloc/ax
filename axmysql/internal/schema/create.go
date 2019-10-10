package schema

import (
	"database/sql"
	"io/ioutil"
	"regexp"
	"strings"
)

// Create executes DDL queries from the given SQL file.
//
// It uses a very naive regexp pattern to identify 'CREATE TABLE' queries in
// order to first drop the table if it already exists.
//
// If you find yourself expanding the behavior of this function its probably
// time to stop and look for a proper schema management solution.
func Create(db *sql.DB, filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	str := string(buf)

	for _, q := range strings.Split(str, ";") {
		if strings.TrimSpace(q) == "" {
			continue
		}

		if m := createTablePattern.FindStringSubmatch(q); m != nil {
			if _, err := db.Exec(`DROP TABLE IF EXISTS ` + m[1]); err != nil {
				return err
			}
		}

		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

var createTablePattern = regexp.MustCompile(`(?i)CREATE\s+TABLE.*?([A-Z_]+)\s+\(`)
