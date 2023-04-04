package property

// The `AppProperties` interface must implement `RootFieldTag`, which is a
// tag key used in Go struct for the top-level params of the sodas application.
// It is used from environ.Retrieve().
// If inner field is a type of struct, inner fields of the struct must tag `json` tag.
type AppProperties interface {
	RootFieldTag() string
}
