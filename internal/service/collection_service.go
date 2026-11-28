package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"board-game-tracker/internal/models"
)

var (
	ErrNotFound   = errors.New("record not found")
	ErrForbidden  = errors.New("no access to this entry")
	ErrValidation = errors.New("validation error")
)

type CollectionRepository interface {
	Create(ctx context.Context, item *models.CollectionItem) error
	GetByID(ctx context.Context, itemID string) (*models.CollectionItem, error)
	List(ctx context.Context, filter *models.CollectionItem) ([]*models.CollectionItem, error)
	Update(ctx context.Context, item *models.CollectionItem) error
	Delete(ctx context.Context, itemID string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type GameGetter interface {
	GetGameByID(ctx context.Context, id string) (*models.Game, error)
}

type CollectionService struct {
	repo     CollectionRepository
	gameRepo GameGetter
}

func NewCollectionService(repo CollectionRepository, gameRepo GameGetter) *CollectionService {
	return &CollectionService{
		repo:     repo,
		gameRepo: gameRepo,
	}
}

func (s *CollectionService) AddGameInCollection(
	ctx context.Context,
	item *models.CollectionItem,
) (*models.CollectionItem, error) {
	if item == nil {
		return nil, fmt.Errorf("%w: empty object", ErrValidation)
	}
	if item.UserID == "" || item.GameID == "" {
		return nil, fmt.Errorf("%w: user_id and game_id required", ErrValidation)
	}

	if _, err := s.gameRepo.GetGameByID(ctx, item.GameID); err != nil {
		return nil, fmt.Errorf("%w: game with id %s not found", ErrNotFound, item.GameID)
	}

	if item.Status == "" {
		item.Status = models.StatusOwned
	} else if !item.Status.IsValid() {
		return nil, fmt.Errorf("%w: incorrect status %s", ErrValidation, item.Status)
	}

	if item.Condition == "" {
		item.Condition = models.ConditionNew
	} else if !item.Condition.IsValid() {
		return nil, fmt.Errorf("%w: incorrect state %s", ErrValidation, item.Condition)
	}

	if item.PurchasePrice != nil && *item.PurchasePrice < 0 {
		return nil, fmt.Errorf("%w: the purchase price cannot be negative", ErrValidation)
	}

	now := time.Now()
	item.ID = uuid.NewString()
	item.CreatedAt = now
	item.UpdatedAt = now

	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *CollectionService) GetAllGamesFromCollection(
	ctx context.Context,
	filter *models.CollectionItem,
) ([]*models.CollectionItem, error) {
	if filter == nil || filter.UserID == "" {
		return nil, fmt.Errorf("%w: user_id not specified", ErrValidation)
	}

	items, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	if items == nil {
		items = make([]*models.CollectionItem, 0)
	}

	return items, nil
}

func (s *CollectionService) UpdateGameInCollection(
	ctx context.Context,
	patch *models.CollectionItem,
) (*models.CollectionItem, error) {
	if patch == nil || patch.ID == "" {
		return nil, fmt.Errorf("%w: post id not specified", ErrValidation)
	}

	existing, err := s.repo.GetByID(ctx, patch.ID)
	if err != nil {
		return nil, ErrNotFound
	}

	if existing.UserID != patch.UserID {
		return nil, ErrForbidden
	}

	if patch.Status != "" {
		if !patch.Status.IsValid() {
			return nil, fmt.Errorf("%w: incorrect status %s", ErrValidation, patch.Status)
		}
		existing.Status = patch.Status
	}

	if patch.Condition != "" {
		if !patch.Condition.IsValid() {
			return nil, fmt.Errorf("%w: incorrect state %s", ErrValidation, patch.Condition)
		}
		existing.Condition = patch.Condition
	}

	if patch.PurchasePrice != nil {
		if *patch.PurchasePrice < 0 {
			return nil, fmt.Errorf("%w: the purchase price cannot be negative", ErrValidation)
		}
		existing.PurchasePrice = patch.PurchasePrice
	}

	if patch.SoldPrice != nil {
		if *patch.SoldPrice < 0 {
			return nil, fmt.Errorf("%w: the selling price cannot be negative", ErrValidation)
		}
		existing.SoldPrice = patch.SoldPrice
	}

	if patch.Notes != "" {
		existing.Notes = patch.Notes
	}

	existing.UpdatedAt = time.Now()

	if err = s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *CollectionService) DeleteGameFromCollection(ctx context.Context, userID, itemID string) error {
	existing, err := s.repo.GetByID(ctx, itemID)
	if err != nil {
		return ErrNotFound
	}
	if existing.UserID != userID {
		return ErrForbidden
	}

	return s.repo.Delete(ctx, itemID)
}

func (s *CollectionService) DeleteCollection(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("%w: user_id not specified", ErrValidation)
	}
	return s.repo.DeleteByUserID(ctx, userID)
}
