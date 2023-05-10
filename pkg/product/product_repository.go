package product

import (
	"errors"
	"gin-exercise/pkg/util"
	"github.com/rs/xid"
	"sync"
	"time"
)

var (
	IdTakenError = "id already registered"
)

type RepositoryIface interface {
	Store(*string, float32, string) (*Product, error)
	Find(string) (*Product, bool)
	Delete(string) bool
	GetMany(int) ([]Product, error)
}

type Repository struct {
	store Store[string, Product]
	ids   util.Set
	lock  sync.Mutex
}

func NewProductRepository() *Repository {
	return &Repository{
		store: NewProductMemStore(),
		ids:   make(map[string]struct{}, 10),
	}
}

func (r *Repository) Store(id *string, price float32, description string) (*Product, error) {
	r.acquireRepoLock()
	if id == nil {
		newId := generateNewId(r)
		return saveInto(r, newId, price, description)
	}

	if r.ids.Present(*id) {
		return nil, errors.New(IdTakenError)
	}
	return saveInto(r, *id, price, description)
}

func (r *Repository) Find(id string) (*Product, bool) {
	r.acquireRepoLock()
	return r.store.Find(id)
}

func (r *Repository) Delete(id string) bool {
	r.acquireRepoLock()
	return r.store.Delete(id)
}

func (r *Repository) GetMany(amount int) ([]Product, error) {
	r.acquireRepoLock()
	maxItems := amount
	if amount > len(r.ids) {
		maxItems = len(r.ids)
	}
	var result []Product

	counter := 0
	for id, _ := range r.ids {
		if counter >= maxItems {
			break
		}
		prd, _ := r.Find(id)
		result = append(result, *prd)
		counter++
	}

	return result, nil
}

func generateNewId(r *Repository) string {
	for {
		candidateId := xid.New().String()
		if r.ids.Present(candidateId) == false {
			return candidateId
		}
	}
}

func saveInto(r *Repository, id string, price float32, description string) (*Product, error) {
	newProduct, newErr := NewProduct(id, price, description, time.Now())
	if newErr != nil {
		return nil, newErr
	}

	savedProduct, storeErr := r.store.Store(id, *newProduct)
	if storeErr != nil {
		return nil, storeErr
	}

	r.ids.Put(id)
	return savedProduct, nil
}

func (r *Repository) acquireRepoLock() {
	r.lock.Lock()
	defer r.lock.Unlock()
}
