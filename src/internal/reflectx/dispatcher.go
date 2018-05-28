package reflectx

import (
	"fmt"
	"reflect"
)

// MakeDispatcher generates a function that performs a type-switch on one of
// it's arguments in order to dispatch the call to a method that handles values
// of that specific type.
//
// fn is a pointer to a function variable that describes the signature to match.
// It is assigned the generated function.
//
// It returns the types that that are supported by the generated type-switch.
func MakeDispatcher(fn interface{}, arg, impl reflect.Type) []reflect.Type {
	if arg.Kind() != reflect.Interface {
		panic(fmt.Sprintf("can not dispatch based on %s, expected an interface", arg))
	}

	v := reflect.ValueOf(fn).Elem()

	d, types := makeDispatcher(v.Type(), arg, impl)
	v.Set(d)

	return types
}

// makeDispatcher builds the function used by MakeDispatcher.
func makeDispatcher(sig, arg, impl reflect.Type) (reflect.Value, []reflect.Type) {
	pos, err := findInputArg(sig, arg)
	if err != nil {
		panic(err)
	}

	var types []reflect.Type
	methods := map[reflect.Type]reflect.Method{}

	for i := 0; i < impl.NumMethod(); i++ {
		m := impl.Method(i)

		if !matchesSignature(pos, sig, m.Type) {
			continue
		}

		concrete := m.Type.In(pos)
		if _, ok := methods[concrete]; ok {
			panic(fmt.Sprintf(
				"found multiple target methods for %s on %s",
				concrete,
				impl,
			))
		}

		methods[concrete] = m
		types = append(types, concrete)
	}

	fn := reflect.MakeFunc(
		sig,
		func(args []reflect.Value) []reflect.Value {
			t := args[pos].Elem().Type()
			m, ok := methods[t]

			if !ok {
				panic(fmt.Sprintf(
					"could not find target method for %s on %s",
					t,
					impl,
				))
			}

			for i, a := range args {
				args[i] = reflect.ValueOf(
					a.Interface(),
				)
			}

			return m.Func.Call(args)
		},
	)

	return fn, types
}

// matchesSignature returns true if m matches the signature sig.
//
// The receiver of m, and the argument in position pos must be some concrete
// implementation of the interfaces in sig. All other input and output
// parameters must match exactly.
func matchesSignature(pos int, sig, m reflect.Type) bool {
	if sig.NumIn() != m.NumIn() {
		return false
	}

	if sig.NumOut() != m.NumOut() {
		return false
	}

	for i := 0; i < m.NumIn(); i++ {
		sp := sig.In(i)
		mp := m.In(i)

		if i == 0 || i == pos {
			if mp.Kind() == reflect.Interface {
				return false
			}

			if !mp.Implements(sp) {
				return false
			}
		} else if mp != sp {
			return false
		}
	}

	for i := 0; i < m.NumOut(); i++ {
		sp := sig.Out(i)
		mp := m.Out(i)

		if mp != sp {
			return false
		}
	}

	return true
}

// findInputArg returns the index of the first argument to fn which is of type arg.
func findInputArg(fn, arg reflect.Type) (int, error) {
	for i := 1; i < fn.NumIn(); i++ {
		a := fn.In(i)

		if a == arg {
			return i, nil
		}
	}

	return 0, fmt.Errorf(
		"signature %s does not accept an input parameter of type %s",
		fn,
		arg,
	)
}
