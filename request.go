package gin_grpc

import "reflect"

// Request encapsulated the rpc request.
type Request struct {
	req any
}

// Bind binds the rpc request to the given value.
func (r *Request) Bind(req any) {
	if reflect.TypeOf(req).Kind() != reflect.Ptr {
		panic("req must be a pointer")
	}

	reflect.ValueOf(req).Elem().Set(reflect.ValueOf(r.req))

	return
}

// GetField returns the value of the given field from the rpc request.
func (r *Request) GetField(field string) any {
	return reflect.ValueOf(r.req).Elem().FieldByName(field).Interface()
}

func (r *Request) getAllFields() map[string]any {
	m := make(map[string]any)
	v := reflect.ValueOf(r.req).Elem()
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).CanInterface() {
			m[v.Type().Field(i).Name] = v.Field(i).Interface()
		}
	}
	return m
}
