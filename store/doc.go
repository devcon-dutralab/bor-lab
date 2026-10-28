// Package store defines the minimal key-value store the simulation persists
// module state through. It models the shape of the Cosmos SDK KVStore contract
// without depending on the real iavl or cosmos-db trees, which keeps the
// simulation's build fast and its dependency surface small.
package store
