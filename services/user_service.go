package services

import (
	"context"
	"errors"
	"strings"

	"ptm-indonesia/config"
	"ptm-indonesia/model"
	repositoryContract "ptm-indonesia/repository/contract"
	servicesContract "ptm-indonesia/services/contract"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserEmailAlreadyUsed = errors.New("user email already used")
	ErrUserInvalidStatus    = errors.New("user status is invalid")
	ErrUserInvalidSort      = errors.New("user sort is invalid")
)

type UserService struct {
	userRepository repositoryContract.UserRepository
	adminEmail     string
}

func NewUserService(cfg *config.AppConfig, userRepository repositoryContract.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
		adminEmail:     normalizeEmail(cfg.Admin.Email),
	}
}

func (s *UserService) List(ctx context.Context, request model.UserListRequest) (*model.UserListResponse, error) {
	request.Normalize()
	if err := validateUserListSort(request); err != nil {
		return nil, err
	}

	totalItems, err := s.userRepository.Count(ctx, request, s.adminEmail)
	if err != nil {
		return nil, err
	}

	users, err := s.userRepository.List(ctx, request, s.adminEmail)
	if err != nil {
		return nil, err
	}

	items := make([]model.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, toUserResponse(&user))
	}

	return &model.UserListResponse{
		Items: items,
		Meta: model.PaginationMeta{
			Page:       request.Page,
			Limit:      request.Limit,
			TotalItems: totalItems,
			TotalPages: calculateTotalPages(totalItems, request.Limit),
		},
	}, nil
}

func (s *UserService) FindByID(ctx context.Context, id uint64) (*model.UserResponse, error) {
	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *UserService) Create(ctx context.Context, request model.UserCreateRequest, actorID uint64) (*model.UserResponse, error) {
	status, err := resolveValidatedUserStatus(request.Status)
	if err != nil {
		return nil, err
	}

	email := normalizeEmail(request.Email)
	exists, err := s.userRepository.ExistsByEmail(ctx, email, nil)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserEmailAlreadyUsed
	}

	user := &model.User{
		Name:      strings.TrimSpace(request.Name),
		Email:     email,
		Address:   normalizeNullableString(request.Address),
		Telp:      normalizeNullableString(request.Telp),
		Status:    status,
		StatusRow: model.ActiveUserStatusRow(),
		CreatedBy: actorID,
		UpdatedBy: actorID,
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, mapUserPersistenceError(err)
	}

	response := toUserResponse(user)
	return &response, nil
}

func (s *UserService) Update(ctx context.Context, id uint64, request model.UserUpdateRequest, actorID uint64) (*model.UserResponse, error) {
	status, err := resolveValidatedUserStatus(request.Status)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	email := normalizeEmail(request.Email)
	exists, err := s.userRepository.ExistsByEmail(ctx, email, &id)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrUserEmailAlreadyUsed
	}

	user.Name = strings.TrimSpace(request.Name)
	user.Email = email
	user.Address = normalizeNullableString(request.Address)
	user.Telp = normalizeNullableString(request.Telp)
	user.Status = status
	user.UpdatedBy = actorID

	if err := s.userRepository.Update(ctx, user); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, mapUserPersistenceError(err)
	}

	updatedUser, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	response := toUserResponse(updatedUser)
	return &response, nil
}

func (s *UserService) Delete(ctx context.Context, id uint64, actorID uint64) error {
	_, err := s.userRepository.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	if err := s.userRepository.SoftDelete(ctx, id, actorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return nil
}

func toUserResponse(user *model.User) model.UserResponse {
	return model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		GoogleID:  user.GoogleID,
		AvatarURL: user.AvatarURL,
		Address:   user.Address,
		Telp:      user.Telp,
		Status:    user.Status,
		CreatedBy: user.CreatedBy,
		UpdatedBy: user.UpdatedBy,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}

func calculateTotalPages(totalItems int64, limit int) int {
	if totalItems == 0 || limit <= 0 {
		return 0
	}

	return int((totalItems + int64(limit) - 1) / int64(limit))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeNullableString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func resolveValidatedUserStatus(status *int16) (int16, error) {
	if status == nil || !model.IsValidUserStatus(*status) {
		return 0, ErrUserInvalidStatus
	}

	return *status, nil
}

func validateUserListSort(request model.UserListRequest) error {
	if request.SortBy != "" && !model.IsValidUserListSortBy(request.SortBy) {
		return ErrUserInvalidSort
	}

	if request.SortDirection != "" && !model.IsValidUserListSortDirection(request.SortDirection) {
		return ErrUserInvalidSort
	}

	return nil
}

var _ servicesContract.UserService = (*UserService)(nil)
