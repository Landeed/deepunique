# DeepUnique

Package `deepunique` is a Go package expanding the `unique` package to support the semantics of the `reflect.DeepEqual` function, including adding support for noncomparable types like slices and maps.

This is primarily useful for efficiently filtering slices so that none of the elements are DeepEqual to each other.

This is similar to making a "set", but with more useful semantics. For example, `deepunique` will treat two different pointers to the same value as equal.

## Example

```go
alice := "Alice"
otherAlice := "Alice"

alices := []*string{&alice, &otherAlice}

// This will filter out one of the alices.
uniqueAlices := deepunique.Unique(alices)
```

See [example_test.go](example_test.go) for how this compares to `reflect.DeepEqual` and `unique.Make`.

## Advanced Usage

The unique handles can be used directly, but there's some complexity around maintaining the internal `unique` pointers through the serialization needed to support slices. See [example_test.go](example_test.go).

## Limitations

This package does not currently support recursive types and may encounter issues with `Chan`, `UnsafePointer`, or `Invalid` types.
