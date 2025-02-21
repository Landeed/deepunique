// go_api/pkg/deepunique/deepunique_test.go
package deepunique

import (
	"testing"
)

func TestDeduplicatePointerless(t *testing.T) {
	type testStruct struct {
		ID   int
		Name string
	}

	tests := []struct {
		name     string
		input    []testStruct
		expected []testStruct
	}{
		{
			name: "No duplicates",
			input: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expected: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
		},
		{
			name: "With duplicates",
			input: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expected: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeduplicatePointerless(tt.input)
			for i := range result {
				if result[i].ID != tt.expected[i].ID || result[i].Name != tt.expected[i].Name {
					t.Errorf("expected %v, got %v", tt.expected, result)
					break
				}
			}
		})
	}
}

func TestUniquePointerless(t *testing.T) {
	type testStruct struct {
		ID   int
		Name string
	}

	tests := []struct {
		name     string
		input    []testStruct
		expected []testStruct
	}{
		{
			name: "No duplicates",
			input: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expected: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
		},
		{
			name: "With duplicates",
			input: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
			expected: []testStruct{
				{ID: 1, Name: "Alice"},
				{ID: 2, Name: "Bob"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniquePointerless(tt.input)
			for i := range result {
				if result[i].ID != tt.expected[i].ID || result[i].Name != tt.expected[i].Name {
					t.Errorf("expected %v, got %v", tt.expected, result)
					break
				}
			}
		})
	}
}

func TestDeduplicateWithPointers(t *testing.T) {
	// Shows that DeduplicatePointerless does not work with pointers
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
				{ID: 1, Name: &anotherAlice}, // don't want this
				{ID: 2, Name: &bob},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeduplicatePointerless(tt.input)
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

func TestUniqueWithPointers(t *testing.T) {
	// Shows that UniquePointerless does not work well with pointers
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
				{ID: 1, Name: &anotherAlice}, // don't want this
				{ID: 2, Name: &bob},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniquePointerless(tt.input)
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
