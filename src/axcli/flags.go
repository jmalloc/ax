package axcli

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/pflag"
	strcase "github.com/stoewer/go-strcase"
)

// declareFlags creates flags on fs for each of the fields in the struct v.
// v must be a pointer to a struct.
func declareFlags(fs *pflag.FlagSet, v interface{}) error {
	pt := reflect.TypeOf(v)

	if pt.Kind() != reflect.Ptr {
		return errors.New("v must be a pointer-to-struct")
	}

	st := pt.Elem()
	if st.Kind() != reflect.Struct {
		return errors.New("v must be a pointer-to-struct")
	}

	rv := reflect.ValueOf(v).Elem()

	for i := 0; i < rv.NumField(); i++ {
		f := st.Field(i)

		// exclude special protocol buffers fields
		if strings.HasPrefix(f.Name, "XXX_") {
			continue
		}

		p := rv.Field(i).Addr().Interface()
		n := FlagName(f)
		u := FlagUsage(f)

		if ok := declareFlagForField(fs, p, n, u); !ok {
			return fmt.Errorf(
				"could not generate flag for %s.%s (%s)",
				reflect.TypeOf(v),
				f.Name,
				f.Type,
			)
		}
	}

	return nil
}

// FlagName returns the name of the flag to use for f.
func FlagName(f reflect.StructField) string {
	n := strcase.KebabCase(f.Name)

	switch n {
	case "help", "timeout":
		return "x-" + n
	default:
		return n
	}
}

// FlagUsage returns the usage string to use for the flag that populates f.
func FlagUsage(f reflect.StructField) string {
	return fmt.Sprintf("Populates the '%s' field", f.Name)
}

// declareFlagForField declares a flag that sets the value pointed to by p.
func declareFlagForField(fs *pflag.FlagSet, p interface{}, n, u string) bool {
	switch ptr := p.(type) {
	case *uint:
		fs.UintVar(ptr, n, *ptr, u)
	case *uint8:
		fs.Uint8Var(ptr, n, *ptr, u)
	case *uint16:
		fs.Uint16Var(ptr, n, *ptr, u)
	case *uint32:
		fs.Uint32Var(ptr, n, *ptr, u)
	case *uint64:
		fs.Uint64Var(ptr, n, *ptr, u)
	case *int:
		fs.IntVar(ptr, n, *ptr, u)
	case *int8:
		fs.Int8Var(ptr, n, *ptr, u)
	case *int16:
		fs.Int16Var(ptr, n, *ptr, u)
	case *int32:
		fs.Int32Var(ptr, n, *ptr, u)
	case *int64:
		fs.Int64Var(ptr, n, *ptr, u)
	case *float32:
		fs.Float32Var(ptr, n, *ptr, u)
	case *float64:
		fs.Float64Var(ptr, n, *ptr, u)
	case *string:
		fs.StringVar(ptr, n, *ptr, u)
	case *bool:
		fs.BoolVar(ptr, n, *ptr, u)
	case *time.Duration:
		fs.DurationVar(ptr, n, *ptr, u)
	default:
		return false
	}

	return true
}
