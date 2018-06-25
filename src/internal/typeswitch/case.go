package typeswitch

import (
	"reflect"
)

// Case represents a single case within a TypeSwitch.
type Case struct {
	// Method is the method that is invoked when this case is hit.
	Method reflect.Method

	// InputMap maps the index of the input parameters passed to Call() to the
	// index of the method's input parameters. A value of -1 indicates that the
	// method does require that parameter.
	InputMap []int

	// OutputMap maps the index of the method's output parameters to the index of
	// the output parameters returned by Call(). All of the method's output
	// parameters are mapped, however gaps in the returned output parameters are
	// permitted and are represented by a nil pointer.
	OutputMap []int
}

// Call invokes c.Method, mapping the input and output parameters as per
// c.InputMap and OutputMap, respectively.
//
// The first value of in is the receiver, and is always passed as the first
// parameter to the method.
func (c Case) Call(in []interface{}) []interface{} {
	methodIn := make([]reflect.Value, c.Method.Type.NumIn())

	for from, v := range in {
		to := c.InputMap[from]
		if to != -1 {
			methodIn[to] = reflect.ValueOf(v)
		}
	}

	methodOut := c.Method.Func.Call(methodIn)
	out := make([]interface{}, len(c.OutputMap))

	for from, v := range methodOut {
		to := c.OutputMap[from]
		out[to] = v.Interface()
	}

	return out
}
