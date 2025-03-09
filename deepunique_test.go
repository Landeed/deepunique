// go_api/pkg/deepunique/deepunique_test.go
package deepunique

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"unique"
)

func TestComparePointerHandles(t *testing.T) {
	// Shows that handles of pointers don't work like deep equality.

	alice := "Alice"
	anotherAlice := "Alice"

	handle1 := unique.Make(&alice)
	handle2 := unique.Make(&anotherAlice)

	if handle1 == handle2 {
		t.Errorf("expected different handles, got %v", handle1)
	}
	if !reflect.DeepEqual(handle1, handle2) {
		t.Errorf("expected %v deep equal to %v", handle1, handle2)
	}
}

func TestRecursiveHandle(t *testing.T) {
	// shows that recursive comparisons work.
	// We won't include in the initial implementation of deepunique.

	type testStruct struct {
		strct *testStruct
	}

	alice := testStruct{strct: new(testStruct)}
	alice.strct = &alice

	handle := unique.Make(alice)

	if handle.Value().strct.strct.strct != &alice {
		t.Errorf("expected %v, got %v", &alice, handle.Value().strct)
	}
}

func TestFuncType(t *testing.T) {
	// shows that the best comparison for functions is a handle of the reflect value.

	// funcs can't be compared directly.
	alice := func() bool { return false }
	alice2 := func() bool { return false }

	// Values can be compared, but are just pointer comparisons, not deep.
	// That means we can make handles.
	val := reflect.ValueOf(alice)
	val2 := reflect.ValueOf(alice2)
	if val == val2 {
		t.Errorf("expected %v, got %v", val.IsNil(), val2.IsNil())
	}

	// panics: val.Interface() == val2.Interface()

	otherVal := reflect.ValueOf(alice)
	if otherVal != val {
		t.Errorf("expected %v, got %v", val, otherVal)
	}

	// These two functions have the same signature, so they have the same type.
	alice3 := func() bool { return true }
	Type := reflect.TypeOf(alice)
	Type3 := reflect.TypeOf(alice3)
	if Type != Type3 {
		t.Errorf("expected %v, got %v", Type, Type3)
	}

	alice4 := func() {}
	val4 := reflect.ValueOf(alice4)
	if val4.IsNil() {
		t.Errorf("%v should not be nil", val4)
	}
}

func TestValHandles(t *testing.T) {
	// shows that we can make handles of values.
	// See TestDuplicateTypes for handles of values with same underlying type.
	// Unfortunately doesn't work for values of pointer elements (see TestPointerElemHandle)
	alice := reflect.ValueOf(int32(1))
	otherAlice := reflect.ValueOf(int64(1))
	handle := unique.Make(alice)
	otherHandle := unique.Make(otherAlice)
	if handle == otherHandle {
		t.Errorf("expected %v, got %v", handle, otherHandle)
	}
}

func TestTypeHandle(t *testing.T) {
	// Shows that handles of types work in simple conditions
	// See TestDuplicateTypes for more complex conditions
	alice := "Alice"
	handle := unique.Make(alice)
	val := reflect.ValueOf(handle)
	handleType := unique.Make(val.Type())

	bob := "bob"
	handle2 := unique.Make(bob)
	val2 := reflect.ValueOf(handle2)
	handleType2 := unique.Make(val2.Type())

	if handleType != handleType2 {
		t.Errorf("expected %v, got %v", handleType, handleType2)
	}

	val3 := reflect.ValueOf("test")
	handleType3 := unique.Make(val3.Type())
	if handleType == handleType3 {
		t.Errorf("expected different handles, got %v", handleType)
	}
}

func evilAlice() any {
	type notString string

	return notString("Alice")
}

func evilAlice2() any {
	type notString string

	return notString("Alice")
}

func TestDuplicateTypes(t *testing.T) {
	// Shows that we can serialize the handle pointer for two different types, even if they
	// have the same value. Pointers will never accidentally clash between two types as long
	// as both handles remain in memory.
	// Shows that for uncomparable types, we need to serialize the handle of the value.Type(),
	// not the string.

	alice := evilAlice()
	alice2 := evilAlice2()
	if alice == alice2 {
		t.Errorf("alice and alice2 both %v", alice)
	}
	if reflect.DeepEqual(alice, alice2) {
		// DeepEqual cares about type, not just underlying value.
		t.Errorf("alice and alice2 are deep equal %v", alice)
	}

	aliceType := reflect.TypeOf(alice)
	alice2Type := reflect.TypeOf(alice2)
	if aliceType == alice2Type {
		t.Errorf("aliceType and alice2Type both %v", aliceType)
	}

	aliceHandle := unique.Make(alice)
	alice2Handle := unique.Make(alice2)
	if aliceHandle == alice2Handle {
		// Handle values are different even though both types are `any`
		t.Errorf("aliceHandle and alice2Handle both %v", aliceHandle)
	}

	aliceHandleType := reflect.TypeOf(aliceHandle)
	alice2HandleType := reflect.TypeOf(alice2Handle)
	if aliceHandleType != alice2HandleType {
		// Handle types are the same for two `any` values
		t.Errorf("expected %v, got %v", aliceHandleType, alice2HandleType)
	}

	val := reflect.ValueOf(alice)
	val2 := reflect.ValueOf(alice2)
	if val == val2 {
		t.Errorf("val and val2 are both %v", val)
	}
	if reflect.DeepEqual(val, val2) {
		// DeepEqual cares about type, even when using ValueOf
		t.Errorf("val and val2 are deepequal %v", val)
	}

	valType := val.Type()
	val2Type := val2.Type()
	if valType == val2Type {
		t.Errorf("valType and val2Type are both %v", valType)
	}
	if valType.String() != val2Type.String() {
		// Type strings are the same despite types being different types
		t.Errorf("valType and val2Type are different %v and %v", valType, val2Type)
	}

	valTypeHandle := unique.Make(valType)
	val2TypeHandle := unique.Make(val2Type)
	if valTypeHandle == val2TypeHandle {
		// Handles of the types are different (correctly)
		t.Errorf("valTypeHandle and val2TypeHandle are both %v", valTypeHandle)
	}

	val3 := reflect.ValueOf(evilAlice())
	val3Type := val3.Type()
	val3TypeHandle := unique.Make(val3Type)
	if val3TypeHandle != valTypeHandle {
		// Handles of types correctly identify same types.
		t.Errorf("val3TypeHandle and valTypeHandle are both %v", val3TypeHandle)
	}
}

func TestArrayHandle(t *testing.T) {
	// Shows that handles can't be serialized directly.
	// See TestSerializeHandle for fix.

	alice := []int{1, 2, 3}
	bob := []int{1, 2, 4}

	aliceHandles := make([]unique.Handle[int], len(alice))
	bobHandles := make([]unique.Handle[int], len(bob))

	for i, v := range alice {
		aliceHandles[i] = unique.Make(v)
	}
	for i, v := range bob {
		bobHandles[i] = unique.Make(v)
	}

	aliceJSON, _ := json.Marshal(aliceHandles)
	bobJSON, _ := json.Marshal(bobHandles)
	if string(aliceJSON) != string(bobJSON) {
		// Json is the same because handles serialize as {}
		t.Errorf("expected %v, got %v", string(aliceJSON), string(bobJSON))
	}
	aliceJSONHandle := unique.Make(string(aliceJSON))
	bobJSONHandle := unique.Make(string(bobJSON))
	if aliceJSONHandle != bobJSONHandle {
		// Handles are unfortunately the same
		t.Errorf("expected %v, got %v", aliceJSONHandle, bobJSONHandle)
	}
}

func TestSerializeHandle(t *testing.T) {
	// Shows that handles can be serialized and still work

	alice := evilAlice()
	serializableHandle := NewSerializableHandle(alice)
	serialized, err := json.Marshal(serializableHandle)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	otherAlice := evilAlice()
	otherSerializableHandle := NewSerializableHandle(otherAlice)
	otherSerialized, err := json.Marshal(otherSerializableHandle)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(serialized) != string(otherSerialized) {
		t.Errorf("expected %v, got %v", string(serialized), string(otherSerialized))
	}

	badAlice := evilAlice2()
	badSerializableHandle := NewSerializableHandle(badAlice)
	badSerialized, err := json.Marshal(badSerializableHandle)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if string(serialized) == string(badSerialized) {
		t.Errorf("expected %v, got %v", string(serialized), string(badSerialized))
	}
}

func TestPointerElemHandle(t *testing.T) {
	// Shows that Elem handles are different, but cast value handles and interfaces are the same.

	alice := "Alice"
	aliceVal := reflect.ValueOf(&alice)

	otherAlice := "Alice"
	otherAliceVal := reflect.ValueOf(&otherAlice)

	aliceElem := aliceVal.Elem()
	otherAliceElem := otherAliceVal.Elem()
	if aliceElem == otherAliceElem {
		// Elems are different for some reason, even though "Kind" is string.
		t.Errorf("expected different values, got %v", aliceElem)
	}

	aliceHandle := unique.Make(aliceElem)
	otherAliceHandle := unique.Make(otherAliceElem)
	if aliceHandle == otherAliceHandle {
		// Really annoying that this makes different handles
		t.Errorf("expected different handles, got %v", aliceHandle)
	}

	aliceString := aliceElem.String()
	otherAliceString := otherAliceElem.String()
	if aliceString != otherAliceString {
		t.Errorf("expected %v, got %v", aliceString, otherAliceString)
	}

	aliceHandle2 := unique.Make(aliceString)
	otherAliceHandle2 := unique.Make(otherAliceString)
	if aliceHandle2 != otherAliceHandle2 {
		// This is how we can deep-compare pointers
		t.Errorf("expected %v, got %v", aliceHandle2, otherAliceHandle2)
	}

	aliceElemInterface := aliceElem.Interface()
	otherAliceElemInterface := otherAliceElem.Interface()
	if aliceElemInterface != otherAliceElemInterface {
		// Interfaces can be compared
		t.Errorf("expected %v, got %v", aliceElemInterface, otherAliceElemInterface)
	}

	aliceHandle3 := unique.Make(aliceElemInterface)
	otherAliceHandle3 := unique.Make(otherAliceElemInterface)
	if aliceHandle3 != otherAliceHandle3 {
		// Handles of interfaces can be compared
		t.Errorf("expected %v, got %v", aliceHandle3, otherAliceHandle3)
	}

	eve := evilAlice()
	ev2 := evilAlice2()
	pointerVal := reflect.ValueOf(&eve)
	pointerVal2 := reflect.ValueOf(&ev2)
	if pointerVal == pointerVal2 {
		t.Errorf("expected different pointers, got %v", pointerVal)
	}
	pointerValHandle := unique.Make(pointerVal.Elem().Interface())
	pointerVal2Handle := unique.Make(pointerVal2.Elem().Interface())
	if pointerValHandle == pointerVal2Handle {
		// Type information isn't lost when using Interface.
		t.Errorf("expected different handles, got %v", pointerValHandle)
	}

	eve3 := evilAlice()
	pointerVal3 := reflect.ValueOf(&eve3)
	pointerValHandle3 := unique.Make(pointerVal3.Elem().Interface())
	if pointerValHandle3 != pointerValHandle {
		// Interface equality does work
		t.Errorf("expected %v, got %v", pointerValHandle, pointerValHandle3)
	}

	/*
		TypedAny1 := TypedAny{
			Type:  NewSerializableHandle(pointerVal.Elem().Type()),
			Value: NewSerializableHandle(pointerVal.Elem().Interface()),
		}
		TypedAny2 := TypedAny{
			Type:  NewSerializableHandle(pointerVal2.Elem().Type()),
			Value: NewSerializableHandle(pointerVal2.Elem().Interface()),
		}
		if TypedAny1 == TypedAny2 {
			t.Errorf("expected %v, got %v", TypedAny1, TypedAny2)
		}
		t.Errorf("TypedAny1: %+v", TypedAny1)
		t.Errorf("TypedAny2: %+v", TypedAny2)
	*/
}

func TestMake(t *testing.T) {
	type testStruct struct {
		IDs  []int
		Name *string
	}

	alice := "Alice"
	bob := "Bob"
	anotherAlice := "Alice"

	handle1, deep1, err := Make(testStruct{IDs: []int{1, 2}, Name: &alice})
	//handle1, deep1, err := Make(&alice)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	handle2, deep2, err := Make(testStruct{IDs: []int{1, 2}, Name: &anotherAlice})
	//handle2, deep2, err := Make(&anotherAlice)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	handle3, deep3, err := Make(testStruct{IDs: []int{2}, Name: &bob})
	//handle3, deep3, err := Make(&bob)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	handle4, deep4, err := Make(testStruct{IDs: []int{2, 1}, Name: &alice})
	//handle4, deep4, err := Make(&alice)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if handle1 != handle2 {
		t.Errorf("expected %v, got %v", handle1, handle2)
	}

	if handle1 == handle3 {
		t.Errorf("handle1 and handle3 should be different, got %v", handle1)
	}

	if handle1 == handle4 {
		t.Errorf("handle1 and handle4 should be different, got %v", handle1)
	}

	_ = deep1
	_ = deep2
	_ = deep3
	_ = deep4
}

func TestUnique(t *testing.T) {
	type testStruct struct {
		ID   int
		Name *string
	}

	alice := "Alice"
	bob := "Bob"
	anotherAlice := "Alice"

	tests := []struct {
		name     string
		input    []testStruct
		expected []testStruct
	}{
		{
			name: "No duplicates with pointers",
			input: []testStruct{
				{ID: 1, Name: &alice},
				{ID: 2, Name: &bob},
			},
			expected: []testStruct{
				{ID: 1, Name: &alice},
				{ID: 2, Name: &bob},
			},
		},
		{
			name: "With duplicates having different pointers to same value",
			input: []testStruct{
				{ID: 1, Name: &alice},
				{ID: 1, Name: &anotherAlice},
				{ID: 2, Name: &bob},
			},
			expected: []testStruct{
				{ID: 1, Name: &alice},
				{ID: 2, Name: &bob},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unique(tt.input)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %v, got %v", len(tt.expected), len(result))
				t.Errorf("expected %v, got %v", tt.expected, result)
				for _, item := range result {
					t.Errorf("item handle: %v", unique.Make(item))
				}
				return
			}
			for i := range result {
				if result[i].ID != tt.expected[i].ID || *result[i].Name != *tt.expected[i].Name {
					t.Errorf("expected %v, got %v", tt.expected, result)
					break
				}
			}
		})
	}
}

func TestUniqueTypedAny(t *testing.T) {
	type testStruct struct {
		ID   int
		Name *any
	}

	eve1 := evilAlice()
	eve2 := evilAlice2()

	tests := []struct {
		name     string
		input    []testStruct
		expected []testStruct
	}{
		{
			name: "With pointers to different types of the same value",
			input: []testStruct{
				{ID: 1, Name: &eve1},
				{ID: 2, Name: &eve2},
			},
			expected: []testStruct{
				{ID: 1, Name: &eve1},
				{ID: 2, Name: &eve2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Unique(tt.input)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %v, got %v", len(tt.expected), len(result))
				t.Errorf("expected %v, got %v", tt.expected, result)
				for _, item := range result {
					t.Errorf("item handle: %v", unique.Make(item))
				}
				return
			}
			for i := range result {
				if result[i].ID != tt.expected[i].ID || *result[i].Name != *tt.expected[i].Name {
					t.Errorf("expected %v, got %v", tt.expected, result)
					break
				}
			}
		})
	}
}

func TestTypeCompare(t *testing.T) {
	// Shows that we need to sort maps by handle pointer string, not just value.
	// We know map keys are comparable, so we don't need full json serialization.

	if 'a' != 97 {
		t.Errorf("Expected 'a' to be equal to 97 in simple comparison")
	}

	handle1 := unique.Make('a')
	handle2 := unique.Make(97)
	str1 := fmt.Sprintf("%v", handle1)
	str2 := fmt.Sprintf("%v", handle2)
	if str1 == str2 {
		t.Errorf("Expected str1 and str2 to be different, got %v", str1)
	}

	// operator not defined on interface
	//deep1 := deepValueMake(reflect.ValueOf('a'))
	//deep2 := deepValueMake(reflect.ValueOf(97))

	// invalid operation: typedAny1 > typedAny2 (operator > not defined on struct)
	/*
		typedAny1 := TypedAny{
			Type:  NewSerializableHandle(reflect.TypeOf('a')),
			Value: handle1,
		}
		typedAny2 := TypedAny{
			Type:  NewSerializableHandle(reflect.TypeOf(97)),
			Value: handle2,
		}
	*/
}

func TestMapDeepEqual(t *testing.T) {
	alice := "Alice"
	otherAlice := "Alice"

	tests := []struct {
		name     string
		map1     map[*string]int
		map2     map[*string]int
		expected bool
	}{
		{
			name: "Maps with keys that are pointers to the same value are not deep equal",
			map1: map[*string]int{
				&alice: 1,
			},
			map2: map[*string]int{
				&otherAlice: 1,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if reflect.DeepEqual(tt.map1, tt.map2) != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, reflect.DeepEqual(tt.map1, tt.map2))
			}
		})
	}
}
