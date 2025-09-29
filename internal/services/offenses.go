package services

import (
	"context"
	"math/big"

	"tipjar/internal/database"
	"tipjar/internal/database/sqlc"
	"tipjar/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type OffenseService struct {
	db *database.DB
}

func NewOffenseService(db *database.DB) *OffenseService {
	return &OffenseService{db: db}
}

func (s *OffenseService) CreateOffense(ctx context.Context, jarID, offenseTypeID, reporterID, offenderID int, notes string, costOverride *float64) (*models.Offense, error) {
	var notesText pgtype.Text
	if notes != "" {
		notesText = pgtype.Text{String: notes, Valid: true}
	}

	var costOverrideNumeric pgtype.Numeric
	if costOverride != nil {
		cents := int64(*costOverride * 100)
		costOverrideNumeric = pgtype.Numeric{
			Int:   big.NewInt(cents),
			Exp:   -2,
			Valid: true,
		}
	}

	params := sqlc.CreateOffenseParams{
		JarID:         int32(jarID),
		OffenseTypeID: int32(offenseTypeID),
		ReporterID:    int32(reporterID),
		OffenderID:    int32(offenderID),
		Notes:         notesText,
		CostOverride:  costOverrideNumeric,
	}

	offense, err := s.db.CreateOffense(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.sqlcOffenseToModel(offense), nil
}

func (s *OffenseService) GetOffenseTypesForJar(ctx context.Context, jarID int) ([]models.OffenseType, error) {
	types, err := s.db.ListOffenseTypesForJar(ctx, int32(jarID))
	if err != nil {
		return nil, err
	}

	result := make([]models.OffenseType, len(types))
	for i, t := range types {
		result[i] = *s.sqlcOffenseTypeToModel(t)
	}

	return result, nil
}

func (s *OffenseService) GetAllOffenseTypesForJar(ctx context.Context, jarID int) ([]models.OffenseType, error) {
	types, err := s.db.ListAllOffenseTypesForJar(ctx, int32(jarID))
	if err != nil {
		return nil, err
	}

	result := make([]models.OffenseType, len(types))
	for i, t := range types {
		result[i] = *s.sqlcOffenseTypeToModel(t)
	}

	return result, nil
}

func (s *OffenseService) GetOffenseType(ctx context.Context, offenseTypeID int) (*models.OffenseType, error) {
	offenseType, err := s.db.GetOffenseType(ctx, int32(offenseTypeID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s.sqlcOffenseTypeToModel(offenseType), nil
}

func (s *OffenseService) CreateOffenseType(ctx context.Context, jarID int, name, description, costType string, costAmount *float64, costAction *string) (*models.OffenseType, error) {
	var descText pgtype.Text
	if description != "" {
		descText = pgtype.Text{String: description, Valid: true}
	}

	var costAmountNumeric pgtype.Numeric
	if costAmount != nil {
		cents := int64(*costAmount * 100)
		costAmountNumeric = pgtype.Numeric{
			Int:   big.NewInt(cents),
			Exp:   -2,
			Valid: true,
		}
	}

	var costActionText pgtype.Text
	if costAction != nil && *costAction != "" {
		costActionText = pgtype.Text{String: *costAction, Valid: true}
	}

	params := sqlc.CreateOffenseTypeParams{
		JarID:       int32(jarID),
		Name:        name,
		Description: descText,
		CostType:    costType,
		CostAmount:  costAmountNumeric,
		CostAction:  costActionText,
	}

	offenseType, err := s.db.CreateOffenseType(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.sqlcOffenseTypeToModel(offenseType), nil
}

func (s *OffenseService) UpdateOffenseType(ctx context.Context, offenseTypeID int, name, description, costType string, costAmount *float64, costAction *string) (*models.OffenseType, error) {
	var descText pgtype.Text
	if description != "" {
		descText = pgtype.Text{String: description, Valid: true}
	}

	var costAmountNumeric pgtype.Numeric
	if costAmount != nil {
		cents := int64(*costAmount * 100)
		costAmountNumeric = pgtype.Numeric{
			Int:   big.NewInt(cents),
			Exp:   -2,
			Valid: true,
		}
	}

	var costActionText pgtype.Text
	if costAction != nil && *costAction != "" {
		costActionText = pgtype.Text{String: *costAction, Valid: true}
	}

	params := sqlc.UpdateOffenseTypeParams{
		ID:          int32(offenseTypeID),
		Name:        name,
		Description: descText,
		CostType:    costType,
		CostAmount:  costAmountNumeric,
		CostAction:  costActionText,
	}

	offenseType, err := s.db.UpdateOffenseType(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.sqlcOffenseTypeToModel(offenseType), nil
}

func (s *OffenseService) DeactivateOffenseType(ctx context.Context, offenseTypeID int) error {
	_, err := s.db.DeactivateOffenseType(ctx, int32(offenseTypeID))
	return err
}

func (s *OffenseService) sqlcOffenseToModel(offense sqlc.Offense) *models.Offense {
	var notes *string
	if offense.Notes.Valid {
		notes = &offense.Notes.String
	}

	var costOverride *float64
	if offense.CostOverride.Valid {
		floatVal, _ := offense.CostOverride.Float64Value()
		if floatVal.Valid {
			costOverride = &floatVal.Float64
		}
	}

	return &models.Offense{
		ID:            int(offense.ID),
		JarID:         int(offense.JarID),
		OffenseTypeID: int(offense.OffenseTypeID),
		ReporterID:    int(offense.ReporterID),
		OffenderID:    int(offense.OffenderID),
		Notes:         notes,
		CostOverride:  costOverride,
		Status:        offense.Status,
		CreatedAt:     offense.CreatedAt.Time,
		UpdatedAt:     offense.UpdatedAt.Time,
	}
}

func (s *OffenseService) sqlcOffenseTypeToModel(offenseType sqlc.OffenseType) *models.OffenseType {
	var description *string
	if offenseType.Description.Valid {
		description = &offenseType.Description.String
	}

	var costAmount *float64
	if offenseType.CostAmount.Valid {
		floatVal, _ := offenseType.CostAmount.Float64Value()
		if floatVal.Valid {
			costAmount = &floatVal.Float64
		}
	}

	var costAction *string
	if offenseType.CostAction.Valid {
		costAction = &offenseType.CostAction.String
	}

	return &models.OffenseType{
		ID:          int(offenseType.ID),
		JarID:       int(offenseType.JarID),
		Name:        offenseType.Name,
		Description: description,
		CostType:    offenseType.CostType,
		CostAmount:  costAmount,
		CostAction:  costAction,
		IsActive:    offenseType.IsActive,
		CreatedAt:   offenseType.CreatedAt.Time,
		UpdatedAt:   offenseType.UpdatedAt.Time,
	}
}