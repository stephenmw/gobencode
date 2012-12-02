// Encodes and decodes bencode format.
// TODO: fix errors
package bencode

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

// An Encoder writes bencoded data to an output stream.
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode writes the bencoded bytes of v to the output stream.
func (e *Encoder) Encode(v interface{}) error {

	// For byte slices, use byte string
	if b, ok := v.([]byte); ok {
		fmt.Fprintf(e.w, "%d:%s", len(b), b)
		return nil
	}

	value := reflect.ValueOf(v)

	// Basic types
	switch value.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(e.w, "i%de", value.Int())
		return nil
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(e.w, "i%de", value.Uint())
		return nil
	case reflect.String:
		s := value.String()
		fmt.Fprintf(e.w, "%d:%s", len(s), s)
		return nil
	}

	// complex types

	return errors.New(fmt.Sprintf("Unsupported type %T", v))
}
