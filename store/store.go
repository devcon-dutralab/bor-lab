package store

// KVStore is the byte-oriented key-value contract modules persist through. It
// is the subset of the Cosmos SDK KVStore the simulation actually uses.
type KVStore interface {
	// Get returns a copy of the value stored at key, or nil when absent.
	Get(key []byte) []byte
	// Has reports whether key is present.
	Has(key []byte) bool
	// Set stores value at key.
	Set(key, value []byte)
	// Delete removes key.
	Delete(key []byte)
}
