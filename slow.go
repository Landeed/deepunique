package deepunique

import "reflect"

func SlowUnique[T any](items []T) []T {
	uniqueItems := make([]T, 0, len(items))
	for _, item := range items {
		flag := false
		for _, other := range uniqueItems {
			if reflect.DeepEqual(item, other) {
				flag = true
				break
			}
		}
		if !flag {
			uniqueItems = append(uniqueItems, item)
		}
	}
	return uniqueItems
}
