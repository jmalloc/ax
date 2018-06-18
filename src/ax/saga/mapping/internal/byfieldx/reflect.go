package byfieldx

import (
	"fmt"
	"reflect"
	"strings"
)

// FieldsToKey returns the concatenation of the string fields on v with the
// given names. It returns an empty string if ANY of the fields are empty.
//
// The concatenated fields are separated by a pipe (|) character. Any pipes or
// backlash (\) characters, they are escaped with a preceding backslash.
func FieldsToKey(v interface{}, names []string) (string, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	r := make([]string, len(names))
	t := rv.Type()

	for i, n := range names {
		f, ok := t.FieldByName(n)
		if !ok {
			return "", fmt.Errorf("%s does not contain a field named %s", t, n)
		}

		if f.Type.Kind() != reflect.String {
			return "", fmt.Errorf("%s.%s is not a string", t, n)
		}

		s := rv.FieldByIndex(f.Index).Interface().(string)
		if s == "" {
			return "", nil
		}

		s = strings.Replace(s, `\`, `\\`, -1) // replace slashes first
		s = strings.Replace(s, `|`, `\|`, -1) // then add slashes to pipes
		r[i] = s
	}

	return strings.Join(r, `|`), nil
}
