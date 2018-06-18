package reflectx

import (
	"fmt"
	"reflect"
)

// StringFields returns the values of a set of string fields on v.
func StringFields(v interface{}, names []string) ([]string, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	r := make([]string, len(names))
	t := rv.Type()

	for i, n := range names {
		f, ok := t.FieldByName(n)
		if !ok {
			return nil, fmt.Errorf("%s does not contain a field named %s", t, n)
		}

		if f.Type.Kind() != reflect.String {
			return nil, fmt.Errorf("%s.%s is not a string", t, n)
		}

		r[i] = rv.FieldByIndex(f.Index).Interface().(string)
	}

	return r, nil
}
