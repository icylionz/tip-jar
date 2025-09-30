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
		result[i] = s.rowToOffenseTypeModel(t.ID, t.JarID, t.Name, t.Description, t.CostAmount, t.CostUnit, t.IsActive, t.CreatedAt, t.UpdatedAt)
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
		result[i] = s.rowToOffenseTypeModel(t.ID, t.JarID, t.Name, t.Description, t.CostAmount, t.CostUnit, t.IsActive, t.CreatedAt, t.UpdatedAt)
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

	model := s.rowToOffenseTypeModel(
		offenseType.ID,
		offenseType.JarID,
		offenseType.Name,
		offenseType.Description,
		offenseType.CostAmount,
		offenseType.CostUnit,
		offenseType.IsActive,
		offenseType.CreatedAt,
		offenseType.UpdatedAt,
	)
	return &model, nil
}

func (s *OffenseService) CreateOffenseType(ctx context.Context, jarID int, name, description string, costAmount *float64, costUnit *string) (*models.OffenseType, error) {
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

	var costUnitText pgtype.Text
	if costUnit != nil && *costUnit != "" {
		costUnitText = pgtype.Text{String: *costUnit, Valid: true}
	}

	params := sqlc.CreateOffenseTypeParams{
		JarID:       int32(jarID),
		Name:        name,
		Description: descText,
		CostAmount:  costAmountNumeric,
		CostUnit:    costUnitText,
	}

	offenseType, err := s.db.CreateOffenseType(ctx, params)
	if err != nil {
		return nil, err
	}

	model := s.rowToOffenseTypeModel(
		offenseType.ID,
		offenseType.JarID,
		offenseType.Name,
		offenseType.Description,
		offenseType.CostAmount,
		offenseType.CostUnit,
		offenseType.IsActive,
		offenseType.CreatedAt,
		offenseType.UpdatedAt,
	)
	return &model, nil
}

func (s *OffenseService) UpdateOffenseType(ctx context.Context, offenseTypeID int, name, description string, costAmount *float64, costUnit *string) (*models.OffenseType, error) {
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

	var costUnitText pgtype.Text
	if costUnit != nil && *costUnit != "" {
		costUnitText = pgtype.Text{String: *costUnit, Valid: true}
	}

	params := sqlc.UpdateOffenseTypeParams{
		ID:          int32(offenseTypeID),
		Name:        name,
		Description: descText,
		CostAmount:  costAmountNumeric,
		CostUnit:    costUnitText,
	}

	offenseType, err := s.db.UpdateOffenseType(ctx, params)
	if err != nil {
		return nil, err
	}

	model := s.rowToOffenseTypeModel(
		offenseType.ID,
		offenseType.JarID,
		offenseType.Name,
		offenseType.Description,
		offenseType.CostAmount,
		offenseType.CostUnit,
		offenseType.IsActive,
		offenseType.CreatedAt,
		offenseType.UpdatedAt,
	)
	return &model, nil
}

func (s *OffenseService) DeactivateOffenseType(ctx context.Context, offenseTypeID int) error {
	_, err := s.db.DeactivateOffenseType(ctx, int32(offenseTypeID))
	return err
}

// Helper function to convert any offense type row to model
func (s *OffenseService) rowToOffenseTypeModel(
	id int32,
	jarID int32,
	name string,
	description pgtype.Text,
	costAmount pgtype.Numeric,
	costUnit pgtype.Text,
	isActive bool,
	createdAt pgtype.Timestamp,
	updatedAt pgtype.Timestamp,
) models.OffenseType {
	var desc *string
	if description.Valid {
		desc = &description.String
	}

	var amount *float64
	if costAmount.Valid {
		floatVal, _ := costAmount.Float64Value()
		if floatVal.Valid {
			amount = &floatVal.Float64
		}
	}

	var unit *string
	if costUnit.Valid {
		unit = &costUnit.String
	}

	return models.OffenseType{
		ID:          int(id),
		JarID:       int(jarID),
		Name:        name,
		Description: desc,
		CostAmount:  amount,
		CostUnit:    unit,
		IsActive:    isActive,
		CreatedAt:   createdAt.Time,
		UpdatedAt:   updatedAt.Time,
	}
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