package sets

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](vs ...T) Set[T] {
	s := make(Set[T], len(vs))
	for _, v := range vs {
		s.Insert(v)
	}
	return s
}

func (s Set[T]) Contains(v T) bool {
	return s.ContainsAny(v)
}

func (s Set[T]) ContainsAny(vs ...T) bool {
	for _, v := range vs {
		_, exists := s[v]
		if exists {
			return true
		}
	}
	return false
}

func (s Set[T]) ContainsAll(vs ...T) bool {
	for _, v := range vs {
		_, exists := s[v]
		if !exists {
			return false
		}
	}
	return true
}

func (s Set[T]) Insert(vs ...T) {
	for _, v := range vs {
		s[v] = struct{}{}
	}
}

func (s Set[T]) Delete(vs ...T) {
	for _, v := range vs {
		delete(s, v)
	}
}

func (s Set[T]) Union(other Set[T]) {
	for k := range other {
		s.Insert(k)
	}
}

func (s Set[T]) Intersect(other Set[T]) {
	var toDelete []T
	for k := range s {
		if _, exists := other[k]; !exists {
			toDelete = append(toDelete, k)
		}
	}
	for _, k := range toDelete {
		delete(s, k)
	}
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.AsSlice())
}

func (s *Set[T]) UnmarshalJSON(text []byte) error {
	var slice []T
	if err := json.Unmarshal(text, &slice); err != nil {
		return errors.Wrapf(err, "failed to unmarshal stringset")
	}
	*s = make(Set[T], len(slice))
	for _, elt := range slice {
		(*s)[elt] = struct{}{}
	}
	return nil
}

func (s Set[T]) Clone() Set[T] {
	return maps.Clone(s)
}

// AsSlice returns the set as a slice in a nondeterministic order.
func (s Set[T]) AsSlice() []T {
	rv := make([]T, 0, len(s))
	for x, _ := range s {
		rv = append(rv, x)
	}
	return rv
}

// Creates a new set from the intersection of sets.
func Intersect[T comparable](sets ...Set[T]) Set[T] {
	if len(sets) == 0 {
		return Set[T]{}
	}

	// Sort by set length.  Starting with the smallest set reduces
	// the work we need to do.
	sort.Slice(sets, func(i, j int) bool {
		return len(sets[i]) < len(sets[j])
	})

	base := sets[0].Clone()
	for _, next := range sets[1:] {
		base.Intersect(next)
	}

	return base
}