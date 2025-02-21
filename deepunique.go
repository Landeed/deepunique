package deepunique

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unique"
)

type SerializableHandle[T comparable] struct {
	handle unique.Handle[T]
	Value  string // stringified unique.Handle
}

func NewSerializableHandle[T comparable](value T) SerializableHandle[T] {
	handle := unique.Make(value)
	return SerializableHandle[T]{
		handle: handle,
		//Value:  fmt.Sprintf("%v", reflect.ValueOf(handle).FieldByName("value").UnsafePointer()),
		Value: fmt.Sprintf("%v", handle),
	}
}

type TypedAny struct {
	Type  SerializableHandle[reflect.Type]
	Value any
}

// TODO: return a more consistent type
func deepValueMake(value reflect.Value) any {
	switch value.Kind() {
	case reflect.Array:
		items := make([]any, value.Len())
		for i := 0; i < value.Len(); i++ {
			items[i] = deepValueMake(value.Index(i))
		}
		return TypedAny{
			Type:  NewSerializableHandle(value.Type()),
			Value: items,
		}
	case reflect.Slice:
		items := make([]any, value.Len())
		for i := 0; i < value.Len(); i++ {
			items[i] = deepValueMake(value.Index(i))
		}
		return TypedAny{
			Type:  NewSerializableHandle(value.Type()),
			Value: items,
		}
	case reflect.Interface:
		return TypedAny{
			Type:  NewSerializableHandle(value.Type()),
			Value: deepValueMake(value.Elem()),
		}
	case reflect.Pointer:
		return TypedAny{
			Type:  NewSerializableHandle(value.Type()),
			Value: deepValueMake(value.Elem()),
		}
	case reflect.Struct:
		items := make([]any, value.NumField())
		for i, n := 0, value.NumField(); i < n; i++ {
			items[i] = deepValueMake(value.Field(i))
		}
		return TypedAny{
			Type:  NewSerializableHandle(value.Type()),
			Value: items,
		}
	case reflect.Map:
		items := make([][2]any, value.Len())
		iter := value.MapRange()
		for iter.Next() {
			// I think map keys have to be comparable, but using deepValueMake to be safe.
			key := deepValueMake(iter.Key())
			val := deepValueMake(iter.Value())
			items = append(items, [2]any{key, val})
		}
		return TypedAny{
			Type:  NewSerializableHandle(value.Type()),
			Value: items,
		}
	case reflect.Func:
		if value.IsNil() {
			return TypedAny{
				Type:  NewSerializableHandle(value.Type()),
				Value: nil,
			}
		}
		// Slightly different behavior from reflect.DeepEqual:
		// Performs a pointer comparison instead of always being unique.
		return NewSerializableHandle(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// A simple unique.Make(value) fails when value is an Elem of a pointer.
		// Have to use TypedAny with the cast value.
		return NewSerializableHandle(value.Interface())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return NewSerializableHandle(value.Interface())
	case reflect.String:
		return NewSerializableHandle(value.Interface())
	case reflect.Bool:
		return NewSerializableHandle(value.Interface())
	case reflect.Float32, reflect.Float64:
		return NewSerializableHandle(value.Interface())
	case reflect.Complex64, reflect.Complex128:
		return NewSerializableHandle(value.Interface())
	case reflect.Chan, reflect.UnsafePointer, reflect.Invalid:
		// Not sure what reflect.DeepEqual is doing here.
		// This might work.
		return NewSerializableHandle(value.Interface())
	default:
		// unreachable with current reflect version
		return NewSerializableHandle(value.Interface())
	}
}

func Make[T any](value T) (unique.Handle[string], any, error) {
	// TODO: simplify holding the pointers in deep.
	deep := deepValueMake(reflect.ValueOf(value))
	// return unique.Make(deep) // Compiles, but panics
	serialized, err := json.Marshal(deep)
	if err != nil {
		return unique.Handle[string]{}, deep, err
	}
	return unique.Make(string(serialized)), deep, nil
}

func Unique[T any](items []T) ([]T, error) {
	seen := make(map[unique.Handle[string]]struct{})
	deeps := make([]any, 0, len(items))
	result := make([]T, 0, len(items))

	for _, item := range items {
		handle, deep, err := Make(item)
		if err != nil {
			return nil, err
		}
		deeps = append(deeps, deep)
		if _, exists := seen[handle]; !exists {
			seen[handle] = struct{}{}
			result = append(result, item)
		}
	}
	_ = deeps // trick to avoid garbage collection. Not sure it works.
	return result, nil
}
