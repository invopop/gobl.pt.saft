// Package atws provides a client for the AT Web Services.
package atws

// Nullable is a helper method that returns a pointer to the given value, i.e.,
// a nullable version of the value. Useful to handle optional fields of
// primitive types in XML that are defined as pointers in the struct (e.g. *int,
// *bool, *string).
func Nullable[T any](v T) *T {
	return &v
}
