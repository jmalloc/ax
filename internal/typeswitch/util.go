package typeswitch

import (
	"reflect"
	"strings"
)

func typesToString(types []reflect.Type) string {
	s := make([]string, len(types))

	for i, t := range types {
		s[i] = t.String()
	}

	return strings.Join(s, ", ")
}

func indexOf(t reflect.Type, types []reflect.Type) int {
	for i, v := range types {
		if t == v {
			return i
		}
	}

	return -1
}
