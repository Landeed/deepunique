package deepunique

import "unique"

func DeduplicatePointerless[T comparable](items []T) []T {
	// tech debt: simplify when go 1.23 is supported
	seen := make(map[T]struct{})
	result := make([]T, 0, len(items))

	for _, item := range items {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func UniquePointerless[T comparable](items []T) []T {
	// Create a map to store unique handles
	seen := make(map[unique.Handle[T]]struct{})
	result := make([]T, 0, len(items))

	for _, item := range items {
		// Create a unique handle for the item
		handle := unique.Make(item)
		if _, exists := seen[handle]; !exists {
			seen[handle] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
