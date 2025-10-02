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

func (s *OffenseService) SetOffenseTypeActiveStatus(ctx context.Context, offenseTypeID int, isActive bool) error {
	_, err := s.db.SetOffenseTypeActiveStatus(ctx, sqlc.SetOffenseTypeActiveStatusParams{
		ID:       int32(offenseTypeID),
		IsActive: isActive,
	})
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
func (s *OffenseService) GetOffenseDetail(ctx context.Context, offenseID int) (*models.OffenseDetail, error) {
	
	offense, err := s.db.GetOffense(ctx, int32(offenseID))
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	offenseType, err := s.db.GetOffenseType(ctx, offense.OffenseTypeID)
	if err != nil {
		return nil, err
	}

	offender, err := s.db.GetUserByID(ctx, offense.OffenderID)
	if err != nil {
		return nil, err
	}

	reporter, err := s.db.GetUserByID(ctx, offense.ReporterID)
	if err != nil {
		return nil, err
	}

	// Calculate amount and unit
	var amount float64
	var unit string

	if offense.CostOverride.Valid {
		floatVal, _ := offense.CostOverride.Float64Value()
		if floatVal.Valid {
			amount = floatVal.Float64
		}
	} else if offenseType.CostAmount.Valid {
		floatVal, _ := offenseType.CostAmount.Float64Value()
		if floatVal.Valid {
			amount = floatVal.Float64
		}
	}

	if offenseType.CostUnit.Valid {
		unit = offenseType.CostUnit.String
	} else {
		unit = "items"
	}

	var notes *string
	if offense.Notes.Valid {
		notes = &offense.Notes.String
	}

	return &models.OffenseDetail{
		ID:              int(offense.ID),
		JarID:           int(offense.JarID),
		OffenseTypeName: offenseType.Name,
		ReporterID:      int(offense.ReporterID),
		ReporterName:    reporter.Name,
		OffenderID:      int(offense.OffenderID),
		OffenderName:    offender.Name,
		Notes:           notes,
		Amount:          amount,
		Unit:            unit,
		Status:          offense.Status,
		CreatedAt:       offense.CreatedAt.Time,
	}, nil
}

func (s *OffenseService) CreatePayment(ctx context.Context, offenseID, userID int, amount *float64, proofURL *string, notes string) (*models.Payment, error) {
	var amountNumeric pgtype.Numeric
	if amount != nil {
		cents := int64(*amount * 100)
		amountNumeric = pgtype.Numeric{
			Int:   big.NewInt(cents),
			Exp:   -2,
			Valid: true,
		}
	}

	var proofURLText pgtype.Text
	if proofURL != nil && *proofURL != "" {
		proofURLText = pgtype.Text{String: *proofURL, Valid: true}
	}

	var proofType pgtype.Text
	if proofURL != nil {
		proofType = pgtype.Text{String: "image", Valid: true}
	}

	params := sqlc.CreatePaymentParams{
		OffenseID: int32(offenseID),
		UserID:    int32(userID),
		Amount:    amountNumeric,
		ProofType: proofType,
		ProofUrl:  proofURLText,
	}

	payment, err := s.db.CreatePayment(ctx, params)
	if err != nil {
		return nil, err
	}

	return s.sqlcPaymentToModel(payment), nil
}

func (s *OffenseService) UpdateOffenseStatus(ctx context.Context, offenseID int, status string) error {
	_, err := s.db.UpdateOffenseStatus(ctx, sqlc.UpdateOffenseStatusParams{
		ID:     int32(offenseID),
		Status: status,
	})
	return err
}

func (s *OffenseService) sqlcPaymentToModel(payment sqlc.Payment) *models.Payment {
	var amount *float64
	if payment.Amount.Valid {
		floatVal, _ := payment.Amount.Float64Value()
		if floatVal.Valid {
			amount = &floatVal.Float64
		}
	}

	var proofType *string
	if payment.ProofType.Valid {
		proofType = &payment.ProofType.String
	}

	var proofURL *string
	if payment.ProofUrl.Valid {
		proofURL = &payment.ProofUrl.String
	}

	var verifiedBy *int
	if payment.VerifiedBy.Valid {
		verifiedByInt := int(payment.VerifiedBy.Int32)
		verifiedBy = &verifiedByInt
	}

	return &models.Payment{
		ID:         int(payment.ID),
		OffenseID:  int(payment.OffenseID),
		UserID:     int(payment.UserID),
		Amount:     amount,
		ProofType:  proofType,
		ProofURL:   proofURL,
		Verified:   payment.Verified,
		VerifiedBy: verifiedBy,
		CreatedAt:  payment.CreatedAt.Time,
		UpdatedAt:  payment.UpdatedAt.Time,
	}
}