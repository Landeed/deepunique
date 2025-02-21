// go_api/pkg/deepunique/deepunique_test.go
package deepunique

import (
	"testing"
	"unique"
)

func TestSlowUnique(t *testing.T) {
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
			result := SlowUnique(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected length %v, got %v", len(tt.expected), len(result))
				t.Errorf("expected %v, got %v", tt.expected, result)
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

func TestSlowUniqueTypedAny(t *testing.T) {
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
			result := SlowUnique(tt.input)
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
