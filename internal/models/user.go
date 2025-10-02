package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Avatar    *string   `json:"avatar" db:"avatar"`
	GoogleID  string    `json:"google_id" db:"google_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type TipJar struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	InviteCode  string    `json:"invite_code" db:"invite_code"`
	CreatedBy   int       `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
type DashboardJar struct {
	*TipJar
	MemberCount int `json:"member_count"`
}

type JarMembership struct {
	ID       int       `json:"id" db:"id"`
	JarID    int       `json:"jar_id" db:"jar_id"`
	UserID   int       `json:"user_id" db:"user_id"`
	Role     string    `json:"role" db:"role"` // 'admin' or 'member'
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}
type JarMemberInfo struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}

type OffenseType struct {
	ID          int       `json:"id" db:"id"`
	JarID       int       `json:"jar_id" db:"jar_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	CostAmount  *float64  `json:"cost_amount" db:"cost_amount"` 
	CostUnit    *string   `json:"cost_unit" db:"cost_unit"` 
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Offense struct {
	ID            int       `json:"id" db:"id"`
	JarID         int       `json:"jar_id" db:"jar_id"`
	OffenseTypeID int       `json:"offense_type_id" db:"offense_type_id"`
	ReporterID    int       `json:"reporter_id" db:"reporter_id"`
	OffenderID    int       `json:"offender_id" db:"offender_id"`
	Notes         *string   `json:"notes" db:"notes"`
	CostOverride  *float64  `json:"cost_override" db:"cost_override"`
	Status        string    `json:"status" db:"status"` // 'pending', 'paid', 'disputed', 'forgiven'
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type Payment struct {
	ID         int       `json:"id" db:"id"`
	OffenseID  int       `json:"offense_id" db:"offense_id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Amount     *float64  `json:"amount" db:"amount"`
	ProofType  *string   `json:"proof_type" db:"proof_type"` // 'image', 'receipt', 'video'
	ProofURL   *string   `json:"proof_url" db:"proof_url"`
	Verified   bool      `json:"verified" db:"verified"`
	VerifiedBy *int      `json:"verified_by" db:"verified_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}
type JarActivity struct {
	ID              int       `json:"id"`
	OffenseTypeName string    `json:"offense_type_name"`
	ReporterID      int       `json:"reporter_id"` 
	ReporterName    string    `json:"reporter_name"`
	ReporterAvatar  *string   `json:"reporter_avatar"`
	OffenderID      int       `json:"offender_id"` 
	OffenderName    string    `json:"offender_name"`
	OffenderAvatar  *string   `json:"offender_avatar"`
	Notes           *string   `json:"notes"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

type MemberBalance struct {
	UserID         int     `json:"user_id"`
	Name           string  `json:"name"`
	Avatar         *string `json:"avatar"`
	TotalOwed      float64 `json:"total_owed"`
	PendingCount   int     `json:"pending_count"`
}

type MemberBalanceByUnit struct {
	UserID       int     `json:"user_id"`
	Name         string  `json:"name"`
	Avatar       *string `json:"avatar"`
	Unit         string  `json:"unit"`
	TotalOwed    float64 `json:"total_owed"`
	OffenseCount int     `json:"offense_count"`
}

type MemberBalanceSummary struct {
	UserID        int                    `json:"user_id"`
	Name          string                 `json:"name"`
	Avatar        *string                `json:"avatar"`
	Balances      []MemberBalanceByUnit  `json:"balances"`
	TotalOffenses int                    `json:"total_offenses"`
}

type OffenseDetail struct {
	ID              int       `json:"id"`
	JarID           int       `json:"jar_id"` // Add this line
	OffenseTypeName string    `json:"offense_type_name"`
	ReporterID      int       `json:"reporter_id"`
	ReporterName    string    `json:"reporter_name"`
	OffenderID      int       `json:"offender_id"`
	OffenderName    string    `json:"offender_name"`
	Notes           *string   `json:"notes"`
	Amount          float64   `json:"amount"`
	Unit            string    `json:"unit"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}