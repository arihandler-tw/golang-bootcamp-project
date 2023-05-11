package inmemory

import "gin-exercise/pkg/product/model"

type Store[K comparable, V any] interface {
	Find(id K) (*V, bool)
	Store(id K, value V) (*V, error)
	Delete(id K) bool
}

type MemStore[K comparable, V any] struct {
	store map[K]V
}

func NewProductMemStore() Store[string, model.Product] {
	return &MemStore[string, model.Product]{store: make(map[string]model.Product)}
}

func (s *MemStore[K, V]) Find(id K) (*V, bool) {
	if value, ok := s.store[id]; ok {
		return &value, true
	}
	return nil, false
}

func (s *MemStore[K, V]) Store(id K, value V) (*V, error) {
	s.store[id] = value
	return &value, nil
}

func (s *MemStore[K, V]) Delete(id K) bool {
	if _, ok := s.store[id]; ok {
		delete(s.store, id)
		return true
	}
	return false
}
