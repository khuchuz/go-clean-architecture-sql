package bookmark

import (
	"context"
	"github.com/khuchuz/go-clean-architecture-sql/models"
)

type Repository interface {
	CreateBookmark(ctx context.Context, user *models.User, bm *models.Bookmark) error
	GetBookmarks(ctx context.Context, user *models.User) ([]*models.Bookmark, error)
	DeleteBookmark(ctx context.Context, user *models.User, id string) error
}
