package ident

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// ErrEmptyID indicates an attempt was made to use an uninitialized ID value.
var ErrEmptyID = errors.New("ID is empty")

// ErrIDNotEmpty indicates an attempt was made to overwrite an existing ID value.
var ErrIDNotEmpty = errors.New("can not replace existing ID")

// ID is a string-based unique identifier.
//
// Many objects use a string as an identifier, and in many cases that string is
// a UUID.
//
// ID encapsulates some common behavior of such IDs, such as simple
// validation and rendering logic without *requiring* that they be a UUID.
//
// IDs are represented as a regular string when encoded as JSON.
//
// ID is intended to be embedded into a more specific identifier type.
type ID struct {
	value string
}

// GenerateUUID sets the ID to a new random UUID.
// It panics if the ID is not empty.
func (id *ID) GenerateUUID() {
	if id.value != "" {
		panic(ErrIDNotEmpty)
	}

	id.value = uuid.NewV4().String()
}

// Parse sets the ID to s. It returns an error if s is empty.
func (id *ID) Parse(s string) error {
	if id.value != "" {
		return ErrIDNotEmpty
	}

	if s == "" {
		return ErrEmptyID
	}

	id.value = strings.ToLower(s)

	return nil
}

// MustParse sets the ID to s. It panics if s is empty.
func (id *ID) MustParse(s string) {
	if err := id.Parse(s); err != nil {
		panic(err)
	}
}

// Get returns the ID as a string, or panics if it is invalid.
func (id ID) Get() string {
	if id.value == "" {
		panic(ErrEmptyID)
	}

	return id.value
}

// String returns a human-readable representation of the ID, which may not be
// the complete or valid ID.
//
// Use Value() to get a copy a valid ID represented as a string.
func (id ID) String() string {
	if id.value == "" {
		return "<unidentified>"
	}

	return FormatID(id.value)
}

// Validate returns an error if the ID is not valid.
func (id ID) Validate() error {
	if id.value == "" {
		return ErrEmptyID
	}

	return nil
}

// MustValidate panics if the ID is not valid.
func (id ID) MustValidate() {
	if err := id.Validate(); err != nil {
		panic(err)
	}
}

// MarshalJSON implements the json.Marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.value)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (id *ID) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &id.value)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id.value), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (id *ID) UnmarshalText(data []byte) error {
	id.value = string(data)
	return nil
}

// Value implements the driver.Valuer interface for database serialization.
func (id ID) Value() (driver.Value, error) {
	if id.value == "" {
		return nil, nil
	}

	return id.value, nil
}

// Scan implements the sql.Scanner interface for database deserialization.
func (id *ID) Scan(v interface{}) error {
	switch x := v.(type) {
	case string:
		id.value = x
	case []byte:
		id.value = string(x)
	case int64:
		id.value = strconv.FormatInt(x, 10)
	case nil:
		id.value = ""
	default:
		return fmt.Errorf("can not scan ID into %s", reflect.TypeOf(v))
	}

	return nil
}
