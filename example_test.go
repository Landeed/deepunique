package deepunique

import (
	"fmt"
	"reflect"
	"unique"
)

func ExampleUnique() {
	alice := "Alice"
	otherAlice := "Alice"

	alices := []*string{&alice, &otherAlice}
	if !reflect.DeepEqual(&alice, &otherAlice) {
		fmt.Println("Expected alice and otherAlice to be deep equal")
	}

	if unique.Make(&alice) == unique.Make(&otherAlice) {
		fmt.Println("Expected alice and otherAlice to have different unique handles")
	}

	uniqueAlices, err := Unique(alices)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Length of uniqueAlices:", len(uniqueAlices))
	// Output: Length of uniqueAlices: 1
}

func ExampleMake() {
	alice := "Alice"
	deepHandle, deepPointers, err := Make(&alice)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	otherAlice := "Alice"
	otherDeepHandle, otherDeepPointers, err := Make(&otherAlice)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Avoid "unused variable" errors. I hope this preserves the pointers through
	// garbage collection.
	_ = deepPointers
	_ = otherDeepPointers

	fmt.Println("alice and otherAlice have the same deep handle:", deepHandle == otherDeepHandle)
	// Output: alice and otherAlice have the same deep handle: true
}
