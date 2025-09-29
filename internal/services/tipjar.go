package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strconv"

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
func (s *TipJarService) GetJarMembers(ctx context.Context, jarID int) ([]models.JarMemberInfo, error) {
	members, err := s.db.ListJarMembers(ctx, int32(jarID))
	if err != nil {
		return nil, err
	}

	result := make([]models.JarMemberInfo, len(members))
	for i, member := range members {
		avatar := ""
		if member.Avatar.Valid {
			avatar = member.Avatar.String
		}

		result[i] = models.JarMemberInfo{
			ID:       int(member.ID),
			UserID:   int(member.UserID),
			Name:     member.Name,
			Email:    member.Email,
			Avatar:   avatar,
			Role:     member.Role,
			JoinedAt: member.JoinedAt.Time,
		}
	}

	return result, nil
}
func (s *TipJarService) GetJarActivity(ctx context.Context, jarID int, limit int) ([]models.JarActivity, error) {
	offenses, err := s.db.ListOffensesForJar(ctx, sqlc.ListOffensesForJarParams{
		JarID:  int32(jarID),
		Limit:  int32(limit),
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	activities := make([]models.JarActivity, len(offenses))
	for i, offense := range offenses {
		var notes *string
		if offense.Notes.Valid {
			notes = &offense.Notes.String
		}

		activities[i] = models.JarActivity{
			ID:              int(offense.ID),
			OffenseTypeName: offense.OffenseTypeName,
			ReporterName:    offense.ReporterName,
			OffenderName:    offense.OffenderName,
			Notes:           notes,
			Status:          offense.Status,
			CreatedAt:       offense.CreatedAt.Time,
		}
	}

	return activities, nil
}

func (s *TipJarService) GetMemberBalances(ctx context.Context, jarID int) ([]models.MemberBalance, error) {
	// Get all jar members
	members, err := s.db.ListJarMembers(ctx, int32(jarID))
	if err != nil {
		return nil, err
	}

	balances := make([]models.MemberBalance, len(members))

	for i, member := range members {
		// Get user's balance in this jar
		balance, err := s.db.GetUserBalanceInJar(ctx, sqlc.GetUserBalanceInJarParams{
			JarID:      int32(jarID),
			OffenderID: member.UserID,
		})
		if err != nil {
			// If error, default to 0
			balance = 0
		}

		// Convert balance to float64
		var totalOwed float64
		if balance != nil {
			if balanceStr, ok := balance.(string); ok {
				if parsed, err := strconv.ParseFloat(balanceStr, 64); err == nil {
					totalOwed = parsed
				}
			} else if balanceFloat, ok := balance.(float64); ok {
				totalOwed = balanceFloat
			}
		}

		// Get pending offense count
		pendingOffenses, err := s.db.ListPendingOffensesForUser(ctx, member.UserID)
		pendingCount := 0
		if err == nil {
			// Count only offenses for this jar
			for _, offense := range pendingOffenses {
				if offense.JarID == int32(jarID) {
					pendingCount++
				}
			}
		}

		var avatar *string
		if member.Avatar.Valid {
			avatar = &member.Avatar.String
		}

		balances[i] = models.MemberBalance{
			UserID:       int(member.UserID),
			Name:         member.Name,
			Avatar:       avatar,
			TotalOwed:    totalOwed,
			PendingCount: pendingCount,
		}
	}

	return balances, nil
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

	// Create default offense type for new jar
	err = s.createDefaultOffenseType(ctx, jar.ID)
	if err != nil {
		return nil, err
	}

	return s.sqlcTipJarToModel(jar), nil
}

// createDefaultOffenseType creates a default "General Offense" type for a new jar
func (s *TipJarService) createDefaultOffenseType(ctx context.Context, jarID int32) error {
	var descText pgtype.Text
	descText = pgtype.Text{String: "A general offense for any rule breaking", Valid: true}

	var costAmount pgtype.Numeric
	costAmount = pgtype.Numeric{Int: big.NewInt(500), Exp: -2, Valid: true} // $5.00

	params := sqlc.CreateOffenseTypeParams{
		JarID:       jarID,
		Name:        "General Offense",
		Description: descText,
		CostType:    "monetary",
		CostAmount:  costAmount,
	}

	_, err := s.db.CreateOffenseType(ctx, params)
	return err
}
