package reflectx

import (
	"fmt"
	"reflect"
)

// PrettyTypeName returns a human-readable name for the type of v.
func PrettyTypeName(v interface{}) string {
	t := reflect.TypeOf(v)

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return fmt.Sprintf(
		"%s.%s",
		t.PkgPath(),
		t.Name(),
	)
}
