package typeswitch

import (
	"fmt"
	"reflect"
	"strings"
)

// Signature defines the criteria that must be met for methods to be included in
// a message set.
type Signature struct {
	Prefix string
	Suffix string
	In     []reflect.Type
	Out    []reflect.Type
}

// IsMatch returns true if m is a match for s. If it is a match, it returns the
// concrete implementation of sw that m accepts.
func (s *Signature) IsMatch(sw reflect.Type, m reflect.Method) (reflect.Type, bool) {
	if !strings.HasPrefix(m.Name, s.Prefix) {
		return nil, false
	}

	if !strings.HasSuffix(m.Name, s.Suffix) {
		return nil, false
	}

	if len(s.In) != m.Type.NumIn() {
		return nil, false
	}

	if len(s.Out) != m.Type.NumOut() {
		return nil, false
	}

	var ct reflect.Type

	for i, provided := range s.In {
		accepted := m.Type.In(i)

		if i == 0 {
			// ensure that the receiver is an exact match, or an implementation of the
			// provided type.
			if accepted != provided && !accepted.Implements(provided) {
				return nil, false
			}
		} else if provided == sw {
			// ensure that the parameter in the same position as the "case type" is an
			// implementation of that case type.
			if !accepted.Implements(sw) {
				return nil, false
			}

			// store the "case type" to be returned
			ct = accepted
		} else {
			// for all other parameters, the types must be an exact match
			if provided != accepted {
				return nil, false
			}
		}
	}

	return ct, ct != nil
}

// MapInputs generates a mapping between the types in p and s.In.
// Any type present in p, that is not present in s.In is represented by -1.
// It returns an error if there are types in s.In that are not present in p.
func (s *Signature) MapInputs(p []reflect.Type) ([]int, error) {
	m := make([]int, 0, len(p))
	m = append(m, 0) // receiver is always the first parameter
	n := 1

	for _, t := range p[1:] {
		i := indexOf(t, s.In)
		m = append(m, i)

		if i != -1 {
			n++
		}
	}

	if n < len(s.In) {
		return nil, fmt.Errorf(
			"signature %s requires input parameters not provided by %s",
			s,
			typesToString(p),
		)
	}

	return m, nil
}

// MapOutputs generates a mapping between the types in s.Out and p.
// Any type present in s.Out, that is not present in p is represented by -1.
// It returns an error if there are types in p that are not present in s.Out.
func (s *Signature) MapOutputs(p []reflect.Type) ([]int, error) {
	m := make([]int, 0, len(s.Out))

	for _, t := range s.Out {
		i := indexOf(t, p)
		if i != -1 {
			m = append(m, i)
		}
	}

	if len(m) != len(s.Out) {
		return nil, fmt.Errorf(
			"signature %s returns output parameters not expected by %s",
			s,
			typesToString(p),
		)
	}

	return m, nil
}

func (s *Signature) String() string {
	return fmt.Sprintf(
		"%s(%s) (%s)",
		s.Prefix+"***"+s.Suffix,
		typesToString(s.In),
		typesToString(s.Out),
	)
}
