package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"math"
	"strings"
)

// HashKey represents the key of a map.
type HashKey struct {
	Type  Type
	Value uint64
}

// Hashable is satisfied by an object that can be hashed.
type Hashable interface {
	HashKey() HashKey
}

// HashKey gets the hash value of a boolean. It is '1' if true and '0' otherwise.
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{b.Type(), 1}
	}

	return HashKey{b.Type(), 0}
}

// HashKey gets the hash value of an integer.
func (i *Integer) HashKey() HashKey {
	return HashKey{i.Type(), uint64(i.Value)}
}

// HashKey gets the hash value of a float.
func (f *Float) HashKey() HashKey {
	return HashKey{f.Type(), math.Float64bits(f.Value)}
}

// HashKey gets the hash value of a string.
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{s.Type(), h.Sum64()}
}

// HashPair represents a key-value pair for a key and value.
type HashPair struct {
	Key   Object
	Value Object
}

// Hash is a map-like data structure.
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() Type {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("MAP {")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
