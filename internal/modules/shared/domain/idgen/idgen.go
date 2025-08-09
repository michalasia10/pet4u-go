package idgen

// Port defines a unique identifier generator for aggregates.
type Port interface {
    NewID(prefix string) string
}


