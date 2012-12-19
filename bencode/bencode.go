// Encodes bencode format.
package bencode

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

const packageName = "bencode"

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
		_, err := fmt.Fprintf(e.w, "%d:%s", len(b), b)
		return err
	}

	value := reflect.ValueOf(v)

	// Basic types
	switch value.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		_, err := fmt.Fprintf(e.w, "i%de", value.Int())
		return err
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		_, err := fmt.Fprintf(e.w, "i%de", value.Uint())
		return err
	case reflect.String:
		s := value.String()
		_, err := fmt.Fprintf(e.w, "%d:%s", len(s), s)
		return err
	}

	// complex types
	switch value.Type().Kind() {
	case reflect.Slice:
		return e.encodeSlice(value)
	case reflect.Struct:
		return e.encodeStruct(value)
	case reflect.Map:
		return e.encodeMap(value)
	}

	return errors.New(fmt.Sprintf("Unsupported type %T", v))
}

func (e *Encoder) encodeSlice(v reflect.Value) error {
	if _, err := e.w.Write([]byte{'l'}); err != nil {
		return err
	}

	for i := 0; i < v.Len(); i++ {
		if err := e.Encode(v.Index(i).Interface()); err != nil {
			return err
		}
	}

	if _, err := e.w.Write([]byte{'e'}); err != nil {
		return err
	}

	return nil
}

type keyValue struct {
	key   string
	value interface{}
}

type keyValueSlice []keyValue

func (p keyValueSlice) Len() int           { return len(p) }
func (p keyValueSlice) Less(i, j int) bool { return p[i].key < p[j].key }
func (p keyValueSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (e *Encoder) encodeStruct(v reflect.Value) error {
	var keyVals []keyValue

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath == "" { // field is exported
			tag := field.Tag.Get(packageName)
			tagPart := strings.SplitAfter(tag, ",")

			key, value := field.Name, v.Field(i).Interface()
			if tagPart[0] != "" {
				key = tagPart[0]
			}

			keyVals = append(keyVals, keyValue{key, value})
		}
	}

	return e.writeDictionary(keyVals)
}

func (e *Encoder) encodeMap(v reflect.Value) error {
	var keyVals []keyValue

	for _, k := range v.MapKeys() {
		keyVals = append(keyVals, keyValue{k.String(), v.MapIndex(k).Interface()})
	}

	return e.writeDictionary(keyVals)
}

func (e *Encoder) writeDictionary(keyVals []keyValue) error {
	sort.Sort(keyValueSlice(keyVals))

	if _, err := e.w.Write([]byte{'d'}); err != nil {
		return err
	}

	for _, kv := range keyVals {
		if err := e.Encode(kv.key); err != nil {
			return err
		}

		if err := e.Encode(kv.value); err != nil {
			return err
		}
	}

	if _, err := e.w.Write([]byte{'e'}); err != nil {
		return err
	}

	return nil
}
