package gin_grpc

import (
	"fmt"
	"reflect"
)

type Response struct {
	resp any
}

// Set the rpc response.
func (r *Response) Set(resp any) {
	if reflect.TypeOf(r.resp) != reflect.TypeOf(resp) {
		panic(fmt.Errorf("internal error: response type mismatch: %s != %s", reflect.TypeOf(r.resp).Name(), reflect.TypeOf(resp).Name()))
	}
	r.resp = resp
}

// SetField set the rpc response field with the given value
func (r *Response) SetField(field string, value any) {
	v := reflect.ValueOf(r.resp).Elem().FieldByName(field)
	if !v.IsValid() {
		panic(fmt.Errorf("internal error: response field not found: %s", field))
	}
	if v.Type() != reflect.TypeOf(value) {
		panic(fmt.Errorf("internal error: response field type mismatch: %s != %s", v.Type().Name(), reflect.TypeOf(value).Name()))
	}
	v.Set(reflect.ValueOf(value))
}

// SetFields set the rpc response fields with the given values
func (r *Response) SetFields(fields H) {
	for field, value := range fields {
		r.SetField(field, value)
	}
}
