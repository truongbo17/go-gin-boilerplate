package page

// Query describes a bounded list request without transport or storage details.
type Query struct {
	Number int
	Size   int
	Search string
}

// Result contains items and counts; transports choose how to present links.
type Result[T any] struct {
	Items  []T
	Number int
	Size   int
	Total  int64
}
