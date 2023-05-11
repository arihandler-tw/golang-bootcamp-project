package inmemory

import (
	"errors"
	"gin-exercise/pkg/product"
	"gin-exercise/pkg/product/model"
	"gin-exercise/pkg/util"
	"github.com/rs/xid"
	"sync"
	"time"
)

type Repository struct {
	store Store[string, model.Product]
	ids   util.Set
	lock  sync.Mutex
}

func NewProductRepository() *Repository {
	return &Repository{
		store: NewProductMemStore(),
		ids:   make(map[string]struct{}, 10),
	}
}

func (r *Repository) Store(id *string, price float32, description string) (*model.Product, error) {
	r.acquireRepoLock()
	if id == nil {
		newId := generateNewId(r)
		return saveInto(r, newId, price, description)
	}

	if r.ids.Present(*id) {
		return nil, errors.New(product.IdTakenError)
	}
	return saveInto(r, *id, price, description)
}

func (r *Repository) Find(id string) (*model.Product, bool) {
	r.acquireRepoLock()
	return r.store.Find(id)
}

func (r *Repository) Delete(id string) bool {
	r.acquireRepoLock()
	return r.store.Delete(id)
}

func (r *Repository) GetMany(amount int) ([]model.Product, error) {
	r.acquireRepoLock()
	maxItems := amount
	if amount > len(r.ids) {
		maxItems = len(r.ids)
	}
	var result []model.Product

	counter := 0
	for id := range r.ids {
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

func saveInto(r *Repository, id string, price float32, description string) (*model.Product, error) {
	newProduct, newErr := model.NewProduct(id, price, description, time.Now())
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
