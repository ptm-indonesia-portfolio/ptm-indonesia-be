package repository

import (
	"context"
	"strings"
	"time"

	"ptm-indonesia/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.baseQuery(ctx).Where("LOWER(email) = ?", normalizeLookupEmail(email)).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	var user model.User
	if err := r.baseQuery(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) List(ctx context.Context, request model.UserListRequest, excludedEmail string) ([]model.User, error) {
	var users []model.User
	err := r.applyListSorting(r.applyListFilters(r.baseQuery(ctx), request, excludedEmail), request).
		Offset(request.Offset()).
		Limit(request.Limit).
		Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Count(ctx context.Context, request model.UserListRequest, excludedEmail string) (int64, error) {
	var total int64
	if err := r.applyListFilters(r.baseQuery(ctx), request, excludedEmail).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	result := r.baseQuery(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"name":       user.Name,
		"email":      user.Email,
		"address":    user.Address,
		"telp":       user.Telp,
		"status":     user.Status,
		"updated_by": user.UpdatedBy,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) SoftDelete(ctx context.Context, id uint64, updatedBy uint64) error {
	now := time.Now()
	result := r.baseQuery(ctx).Model(&model.User{}).Where("id = ?", id).Updates(map[string]any{
		"deleted_at": now,
		"status_row": gorm.Expr("NULL"),
		"updated_by": updatedBy,
	})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string, excludedID *uint64) (bool, error) {
	query := r.baseQuery(ctx).Where("LOWER(email) = ?", normalizeLookupEmail(email))
	if excludedID != nil {
		query = query.Where("id <> ?", *excludedID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return false, err
	}

	return total > 0, nil
}

func (r *UserRepository) UpdateGoogleIdentity(ctx context.Context, userID uint64, googleID string, avatarURL *string) error {
	updates := map[string]any{
		"google_id":  googleID,
		"avatar_url": avatarURL,
		"updated_by": 0,
	}

	result := r.baseQuery(ctx).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *UserRepository) baseQuery(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Model(&model.User{}).Where("deleted_at IS NULL AND status_row = 1")
}

func (r *UserRepository) applyListFilters(query *gorm.DB, request model.UserListRequest, excludedEmail string) *gorm.DB {
	normalizedExcludedEmail := strings.ToLower(strings.TrimSpace(excludedEmail))
	if normalizedExcludedEmail != "" {
		query = query.Where("LOWER(email) <> ?", normalizedExcludedEmail)
	}

	if request.Search == "" {
		return query
	}

	keyword := "%" + request.Search + "%"

	return query.Where(
		"name ILIKE ? OR email ILIKE ? OR COALESCE(address, '') ILIKE ? OR COALESCE(telp, '') ILIKE ?",
		keyword,
		keyword,
		keyword,
		keyword,
	)
}

func (r *UserRepository) applyListSorting(query *gorm.DB, request model.UserListRequest) *gorm.DB {
	if request.SortBy == "" {
		return query.Order("id DESC")
	}

	direction := model.UserListSortDirectionDesc
	if request.SortDirection != "" {
		direction = request.SortDirection
	}

	return query.
		Order(request.SortBy + " " + direction).
		Order("id DESC")
}

func normalizeLookupEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
