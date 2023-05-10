package util

type Set map[string]struct{}

func (s Set) Put(v string) {
	s[v] = struct{}{}
}

func (s Set) Present(v string) bool {
	_, present := s[v]
	return present
}
