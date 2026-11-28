package models

type CollectionStatus string

const (
	StatusOwned    CollectionStatus = "owned"
	StatusWishlist CollectionStatus = "wishlist"
	StatusPreorder CollectionStatus = "preorder"
	StatusSold     CollectionStatus = "sold"
)

func (s CollectionStatus) IsValid() bool {
	switch s {
	case StatusOwned, StatusWishlist, StatusPreorder, StatusSold:
		return true
	default:
		return false
	}
}
