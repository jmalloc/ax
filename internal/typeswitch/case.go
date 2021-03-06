package typeswitch

import (
	"reflect"
)

// Case represents a single case within a Switch.
type Case struct {
	// Method is the method that is invoked when this case is hit.
	Method reflect.Method

	// InputMap maps the index of the input parameters passed to Call() to the
	// index of the method's input parameters. A value of -1 indicates that the
	// method does not require that parameter.
	InputMap []int

	// OutputMap maps the index of the method's output parameters to the index of
	// the output parameters returned by Call(). A value of -1 indicates that the
	// method produces a value that is not required.
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

	for to, from := range c.OutputMap {
		if from != -1 {
			out[to] = methodOut[from].Interface()
		}
	}

	return out
}
