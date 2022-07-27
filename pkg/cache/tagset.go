package cache

import "sync"

// New TagSet.
func NewTagSet() *TagSet {
	s := TagSet{
		d: make(map[string]struct{}),
	}
	return &s
}

// TagSet is our struct that acts as a set data structure
// with string as members.
type TagSet struct {
	l sync.RWMutex
	d map[string]struct{}
}

// Add method to add a member to the TagSet.
func (s *TagSet) Add(member string) {
	s.l.Lock()
	defer s.l.Unlock()

	s.d[member] = struct{}{}
}

// Remove method to remove a member from the TagSet.
func (s *TagSet) Remove(member string) {
	s.l.Lock()
	defer s.l.Unlock()

	delete(s.d, member)
}

// IsMember method to check if a member is present in the TagSet.
func (s *TagSet) IsMember(member string) bool {
	s.l.RLock()
	defer s.l.RUnlock()

	_, found := s.d[member]
	return found
}

// Members method to retrieve all members of the TagSet.
func (s *TagSet) Members() []string {
	s.l.RLock()
	defer s.l.RUnlock()

	keys := make([]string, 0)
	for k := range s.d {
		keys = append(keys, k)
	}
	return keys
}

// Size method to get the cardinality of the TagSet.
func (s *TagSet) Size() int {
	s.l.RLock()
	defer s.l.RUnlock()

	return len(s.d)
}

// Clear method to remove all members from the TagSet.
func (s *TagSet) Clear() {
	s.l.Lock()
	defer s.l.Unlock()

	s.d = make(map[string]struct{})
}
