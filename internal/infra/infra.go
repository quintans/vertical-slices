package infra

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var ErrDoesNotExist = errors.New("does not exist")
var ErrUniquenessViolation = errors.New("uniqueness violation")

type DB[T any] struct {
	data  map[uuid.UUID]T
	mutex sync.RWMutex
}

func NewDB[T any]() *DB[T] {
	return &DB[T]{
		data: make(map[uuid.UUID]T),
	}
}

func (r *DB[T]) GetByID(id uuid.UUID) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if p, ok := r.data[id]; ok {
		return p, nil
	}

	var zero T
	return zero, ErrDoesNotExist
}

func (r *DB[T]) ListAll() []T {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var list []T
	for _, v := range r.data {
		list = append(list, v)
	}
	return list
}

func (r *DB[T]) Create(id uuid.UUID, p T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.data[id]; ok {
		return ErrUniquenessViolation
	}

	r.data[id] = p
	return nil
}

func (r *DB[T]) Delete(id uuid.UUID) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.data, id)
}

func (r *DB[T]) Update(id uuid.UUID, fn func(p T) (T, error)) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	p, ok := r.data[id]
	if !ok {
		return ErrDoesNotExist
	}

	var err error
	p, err = fn(p)
	if err != nil {
		return err
	}

	r.data[id] = p
	return nil
}
