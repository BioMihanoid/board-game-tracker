package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"board-game-tracker/internal/models"
)

type CollectionRepository struct {
	items map[string]*models.CollectionItem
	rw    sync.RWMutex
}

func NewCollectionRepository() *CollectionRepository {
	return &CollectionRepository{
		items: make(map[string]*models.CollectionItem),
	}
}

func (r *CollectionRepository) Create(_ context.Context, item *models.CollectionItem) error {
	if item == nil {
		return errors.New("CollectionRepository.Create: item is nil")
	}
	if item.ID == "" {
		return errors.New("CollectionRepository.Create: item id is empty")
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	if _, exists := r.items[item.ID]; exists {
		return fmt.Errorf("CollectionRepository.Create: item %s already exists", item.ID)
	}

	clone := cloneCollectionItem(item)
	r.items[clone.ID] = clone
	return nil
}

func (r *CollectionRepository) GetByID(_ context.Context, itemID string) (*models.CollectionItem, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	item, ok := r.items[itemID]
	if !ok {
		return nil, fmt.Errorf("CollectionRepository.GetByID: %w", ErrNotFound)
	}

	return cloneCollectionItem(item), nil
}

func (r *CollectionRepository) List(_ context.Context, filter *models.CollectionItem) ([]*models.CollectionItem, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	res := make([]*models.CollectionItem, 0)
	for _, item := range r.items {
		if filter != nil {
			if filter.UserID != "" && item.UserID != filter.UserID {
				continue
			}
			if filter.GameID != "" && item.GameID != filter.GameID {
				continue
			}
			if filter.Status != "" && item.Status != filter.Status {
				continue
			}
		}

		res = append(res, cloneCollectionItem(item))
	}

	return res, nil
}

func (r *CollectionRepository) Update(_ context.Context, item *models.CollectionItem) error {
	if item == nil {
		return errors.New("CollectionRepository.Update: item is nil")
	}

	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.items[item.ID]; !ok {
		return fmt.Errorf("CollectionRepository.Update: %w", ErrNotFound)
	}

	r.items[item.ID] = cloneCollectionItem(item)
	return nil
}

func (r *CollectionRepository) Delete(_ context.Context, itemID string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.items[itemID]; !ok {
		return fmt.Errorf("CollectionRepository.Delete: %w", ErrNotFound)
	}

	delete(r.items, itemID)
	return nil
}

func (r *CollectionRepository) DeleteByUserID(_ context.Context, userID string) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	for id, item := range r.items {
		if item.UserID == userID {
			delete(r.items, id)
		}
	}

	return nil
}

func cloneCollectionItem(src *models.CollectionItem) *models.CollectionItem {
	clone := *src

	if src.PurchasePrice != nil {
		v := *src.PurchasePrice
		clone.PurchasePrice = &v
	}
	if src.SoldPrice != nil {
		v := *src.SoldPrice
		clone.SoldPrice = &v
	}
	if src.PurchasedAt != nil {
		t := *src.PurchasedAt
		clone.PurchasedAt = &t
	}

	return &clone
}
