package services

import (
	"context"

	"tipjar/internal/database"
	"tipjar/internal/database/sqlc"
	"tipjar/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	user, err := s.db.GetUserByGoogleID(ctx, googleID)
	if err != nil {
		if err == pgx.ErrNoRows {
			// User not found - this is expected for new users
			return nil, nil
		}
		// Actual database error
		return nil, err
	}
	
	return s.sqlcUserToModel(user), nil
}

func (s *UserService) CreateUser(ctx context.Context, email, name, avatar, googleID string) (*models.User, error) {
	var avatarText pgtype.Text
	if avatar != "" {
		avatarText = pgtype.Text{String: avatar, Valid: true}
	}

	params := sqlc.CreateUserParams{
		Email:    email,
		Name:     name,
		Avatar:   avatarText,
		GoogleID: googleID,
	}

	user, err := s.db.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.sqlcUserToModel(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.db.GetUserByID(ctx, int32(userID))
	if err != nil {
		if err == pgx.ErrNoRows {
			// User not found
			return nil, nil
		}
		return nil, err
	}

	return s.sqlcUserToModel(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, userID int, name, avatar string) (*models.User, error) {
	var avatarText pgtype.Text
	if avatar != "" {
		avatarText = pgtype.Text{String: avatar, Valid: true}
	}

	params := sqlc.UpdateUserParams{
		ID:     int32(userID),
		Name:   name,
		Avatar: avatarText,
	}

	user, err := s.db.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.sqlcUserToModel(user), nil
}

func (s *UserService) sqlcUserToModel(user sqlc.User) *models.User {
	var avatar *string
	if user.Avatar.Valid {
		avatar = &user.Avatar.String
	}

	return &models.User{
		ID:        int(user.ID),
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    avatar,
		GoogleID:  user.GoogleID,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}
}