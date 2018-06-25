package typeswitch

import (
	"fmt"
	"reflect"
)

// Switch dispatches calls to a method based on the concrete type of some
// interface value.
type Switch map[reflect.Type]Case

// New returns a type-switch that dispatches to methods based on the concrete
// type of one of the parameters.
//
// in and out are the input and output parameters provided and expected by the
// caller, respectively. in[0] must be some concrete receiver type, that is the
// type with the target methods. in[1] must be an interface defining the "switch
// type", which must be an interface.
//
// sigs is a set of method signatures that define which methods of in[0] (the
// receiver) will be considered as targets for each switch "case".
func New(
	in, out []reflect.Type,
	sigs ...*Signature,
) (Switch, map[*Signature][]reflect.Type, error) {
	r := in[0]
	sw := in[1]

	if sw.Kind() != reflect.Interface {
		panic(fmt.Sprintf(
			"can not switch over non-interface %s",
			sw,
		))
	}

	s := Switch{}
	sm := map[*Signature][]reflect.Type{}

	for i := 0; i < r.NumMethod(); i++ {
		m := r.Method(i)

		for _, sig := range sigs {
			ct, ok := sig.IsMatch(sw, m)
			if !ok {
				continue
			}

			if err := s.addCase(ct, m, in, out, sig); err != nil {
				return nil, nil, err
			}

			sm[sig] = append(sm[sig], ct)
		}
	}

	return s, sm, nil
}

// Dispatch invokes the method associated with the type of in[1].
// It panics if there is no case for that type.
func (s Switch) Dispatch(in ...interface{}) []interface{} {
	ct := reflect.TypeOf(in[1])
	if c, ok := s[ct]; ok {
		return c.Call(in)
	}

	panic(fmt.Sprintf(
		"%s has no case for %s",
		reflect.TypeOf(in[0]),
		ct,
	))
}

// Types returns a slice of all of the case types supported by this switch.
func (s Switch) Types() []reflect.Type {
	types := make([]reflect.Type, 0, len(s))
	for t := range s {
		types = append(types, t)
	}

	return types
}

func (s Switch) addCase(
	t reflect.Type,
	m reflect.Method,
	in, out []reflect.Type,
	sig *Signature,
) error {
	if x, ok := s[t]; ok {
		return fmt.Errorf(
			"%s.%s() and %s() both produce a case for %s",
			m.Type.In(0),
			x.Method.Name,
			m.Name,
			t,
		)
	}

	c := Case{
		Method: m,
	}

	var err error
	c.InputMap, err = sig.MapInputs(in)
	if err != nil {
		return err
	}

	c.OutputMap, err = sig.MapOutputs(out)
	if err != nil {
		return err
	}

	s[t] = c

	return nil
}
