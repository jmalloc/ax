package visitor

import (
	"fmt"
	"reflect"
	"strings"
)

// MakeAcceptor generates a function that accepts a value that implements the
// interface type i, and forwards the call to a method on the visitor type v
// that corresponds to that value's concrete type.
//
// fn must be a pointer to a function. Its signature forms the template that
// methods on v must match to be considered the destination for a concrete type.
// It is assigned the value of the generated function.
//
// pre is a prefix that the method name must have in order to be considered a
// match. The prefix may be empty.
func MakeAcceptor(
	fn interface{},
	i, v reflect.Type,
	pre string,
) []reflect.Type {
	if i.Kind() != reflect.Interface {
		panic("i must be an interface")
	}

	p := reflect.ValueOf(fn)
	if p.Kind() != reflect.Ptr {
		panic("fn must be a pointer-to-function")
	}

	fv := p.Elem()
	if fv.Kind() != reflect.Func {
		panic("*fn must be a function")
	}

	sig := fv.Type()
	if sig.NumIn() < 2 {
		panic("*fn must accept at least 2 input parameters")
	}

	rt := sig.In(0)
	if !v.Implements(rt) {
		panic(fmt.Sprintf("v must implement *fn's receiver type %s", rt))
	}

	pos, ok := findInputArg(sig, i)
	if !ok {
		panic(fmt.Sprintf("*fn must have at least one parameter of type %s", i))
	}

	methods := findMatchingMethods(sig, i, v, pos, pre)

	fv.Set(
		makeAcceptor(sig, v, pos, methods),
	)

	types := make([]reflect.Type, 0, len(methods))

	for t := range methods {
		types = append(types, t)
	}

	return types
}

// findMatchingMethods returns the methods on v that match the signature of sig.
//
// A method matches if the input parameter at index ipos is a concrete
// implementation of the interface i, and all other input and output parameters
// are identical to their corresponding parameters in sig, and the method name
// begins with pre.
func findMatchingMethods(
	sig, i, v reflect.Type,
	ipos int,
	pre string,
) map[reflect.Type]reflect.Method {
	methods := map[reflect.Type]reflect.Method{}

	for n := 0; n < v.NumMethod(); n++ {
		m := v.Method(n)

		if !strings.HasPrefix(m.Name, pre) {
			continue
		}

		if !matchesSignature(sig, i, m.Type, ipos) {
			continue
		}

		t := m.Type.In(ipos)
		if _, ok := methods[t]; ok {
			panic(fmt.Sprintf("found multiple target methods for %s on %s", t, v))
		}

		methods[t] = m
	}

	return methods
}

// matchesSignature returns true if m matches the signature sig.
//
// The receiver of m, and the argument in position pos must be some concrete
// implementation of the interfaces in sig. All other input and output
// parameters must match exactly.
func matchesSignature(sig, i, m reflect.Type, ipos int) bool {
	if sig.NumIn() != m.NumIn() {
		return false
	}

	if sig.NumOut() != m.NumOut() {
		return false
	}

	for n := 1; n < m.NumIn(); n++ {
		arg := m.In(n)

		if n == ipos {
			if arg.Kind() == reflect.Interface {
				return false
			}

			if !arg.Implements(i) {
				return false
			}
		} else {
			if arg != sig.In(n) {
				return false
			}
		}
	}

	for n := 0; n < m.NumOut(); n++ {
		if m.Out(n) != sig.Out(n) {
			return false
		}
	}

	return true
}

// makeAcceptor builds the function used by MakeAcceptor.
func makeAcceptor(
	sig, v reflect.Type,
	ipos int,
	methods map[reflect.Type]reflect.Method,
) reflect.Value {
	return reflect.MakeFunc(
		sig,
		func(args []reflect.Value) []reflect.Value {
			t := args[ipos].Elem().Type()
			m, ok := methods[t]

			if !ok {
				panic(fmt.Sprintf("could not find target method for %s on %s", t, v))
			}

			for n, a := range args {
				args[n] = reflect.ValueOf(
					a.Interface(),
				)
			}

			return m.Func.Call(args)
		},
	)
}

// findInputArg returns the index of the first argument to fn which is of type arg.
func findInputArg(fn, arg reflect.Type) (int, bool) {
	for i := 1; i < fn.NumIn(); i++ {
		a := fn.In(i)

		if a == arg {
			return i, true
		}
	}

	return 0, false
}
