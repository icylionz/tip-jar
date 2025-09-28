package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"tipjar/internal/database"
	"tipjar/internal/database/sqlc"
	"tipjar/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type TipJarService struct {
	db *database.DB
}

func NewTipJarService(db *database.DB) *TipJarService {
	return &TipJarService{db: db}
}

func (s *TipJarService) CreateTipJar(ctx context.Context, name, description string, createdBy int) (*models.TipJar, error) {
	inviteCode, err := s.generateInviteCode()
	if err != nil {
		return nil, err
	}

	var descText pgtype.Text
	if description != "" {
		descText = pgtype.Text{String: description, Valid: true}
	}

	params := sqlc.CreateTipJarParams{
		Name:        name,
		Description: descText,
		InviteCode:  inviteCode,
		CreatedBy:   int32(createdBy),
	}

	jar, err := s.db.CreateTipJar(ctx, params)
	if err != nil {
		return nil, err
	}

	// Create admin membership for creator
	_, err = s.db.CreateJarMembership(ctx, sqlc.CreateJarMembershipParams{
		JarID:  jar.ID,
		UserID: int32(createdBy),
		Role:   "admin",
	})
	if err != nil {
		return nil, err
	}

	return s.sqlcTipJarToModel(jar), nil
}

func (s *TipJarService) CreateTipJarWithInviteCode(ctx context.Context, name, description, inviteCode string, createdBy int) (*models.TipJar, error) {
	var descText pgtype.Text
	if description != "" {
		descText = pgtype.Text{String: description, Valid: true}
	}

	params := sqlc.CreateTipJarParams{
		Name:        name,
		Description: descText,
		InviteCode:  inviteCode,
		CreatedBy:   int32(createdBy),
	}

	jar, err := s.db.CreateTipJar(ctx, params)
	if err != nil {
		return nil, err
	}

	// Create admin membership for creator
	_, err = s.db.CreateJarMembership(ctx, sqlc.CreateJarMembershipParams{
		JarID:  jar.ID,
		UserID: int32(createdBy),
		Role:   "admin",
	})
	if err != nil {
		return nil, err
	}

	return s.sqlcTipJarToModel(jar), nil
}

func (s *TipJarService) GetTipJar(ctx context.Context, jarID int) (*models.TipJar, error) {
	jar, err := s.db.GetTipJar(ctx, int32(jarID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s.sqlcTipJarToModel(jar), nil
}

func (s *TipJarService) GetTipJarByInviteCode(ctx context.Context, inviteCode string) (*models.TipJar, error) {
	jar, err := s.db.GetTipJarByInviteCode(ctx, inviteCode)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s.sqlcTipJarToModel(jar), nil
}

func (s *TipJarService) ListTipJarsForUser(ctx context.Context, userID int) ([]*models.TipJar, error) {
	jars, err := s.db.ListTipJarsForUser(ctx, int32(userID))
	if err != nil {
		return nil, err
	}

	result := make([]*models.TipJar, len(jars))
	for i, jar := range jars {
		result[i] = s.sqlcTipJarToModel(sqlc.TipJar{
			ID:          jar.ID,
			Name:        jar.Name,
			Description: jar.Description,
			InviteCode:  jar.InviteCode,
			CreatedBy:   jar.CreatedBy,
			CreatedAt:   jar.CreatedAt,
			UpdatedAt:   jar.UpdatedAt,
		})
	}

	return result, nil
}

func (s *TipJarService) JoinTipJar(ctx context.Context, jarID, userID int) error {
	params := sqlc.CreateJarMembershipParams{
		JarID:  int32(jarID),
		UserID: int32(userID),
		Role:   "member",
	}

	_, err := s.db.CreateJarMembership(ctx, params)
	return err
}

func (s *TipJarService) IsUserJarMember(ctx context.Context, jarID, userID int) (bool, error) {
	return s.db.IsUserJarMember(ctx, sqlc.IsUserJarMemberParams{
		JarID:  int32(jarID),
		UserID: int32(userID),
	})
}

func (s *TipJarService) IsUserJarAdmin(ctx context.Context, jarID, userID int) (bool, error) {
	return s.db.IsUserJarAdmin(ctx, sqlc.IsUserJarAdminParams{
		JarID:  int32(jarID),
		UserID: int32(userID),
	})
}

func (s *TipJarService) generateInviteCode() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:8], nil
}

func (s *TipJarService) sqlcTipJarToModel(jar sqlc.TipJar) *models.TipJar {
	var description *string
	if jar.Description.Valid {
		description = &jar.Description.String
	}

	return &models.TipJar{
		ID:          int(jar.ID),
		Name:        jar.Name,
		Description: description,
		InviteCode:  jar.InviteCode,
		CreatedBy:   int(jar.CreatedBy),
		CreatedAt:   jar.CreatedAt.Time,
		UpdatedAt:   jar.UpdatedAt.Time,
	}
}