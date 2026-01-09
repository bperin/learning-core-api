package utils

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// SQL null type conversion utilities
// These functions help convert between Go types and SQL null types

// String conversions
func SqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func NullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}

func StringToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func NullStringToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// UUID conversions
func SqlNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: *u, Valid: true}
}

// NullUUIDToPtr is already defined in conversions.go

func UUIDToNullUUID(u uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{UUID: u, Valid: true}
}

func NullUUIDToUUID(nu uuid.NullUUID) uuid.UUID {
	if !nu.Valid {
		return uuid.Nil
	}
	return nu.UUID
}

// Int32 conversions
func SqlNullInt32(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func NullInt32ToPtr(ni sql.NullInt32) *int32 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int32
}

func Int32ToNullInt32(i int32) sql.NullInt32 {
	return sql.NullInt32{Int32: i, Valid: true}
}

func NullInt32ToInt32(ni sql.NullInt32) int32 {
	if !ni.Valid {
		return 0
	}
	return ni.Int32
}

// Int64 conversions
func SqlNullInt64(i *int64) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: *i, Valid: true}
}

func NullInt64ToPtr(ni sql.NullInt64) *int64 {
	if !ni.Valid {
		return nil
	}
	return &ni.Int64
}

func Int64ToNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{Int64: i, Valid: true}
}

func NullInt64ToInt64(ni sql.NullInt64) int64 {
	if !ni.Valid {
		return 0
	}
	return ni.Int64
}

// Float64 conversions
func SqlNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{Valid: false}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

func NullFloat64ToPtr(nf sql.NullFloat64) *float64 {
	if !nf.Valid {
		return nil
	}
	return &nf.Float64
}

func Float64ToNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{Float64: f, Valid: true}
}

func NullFloat64ToFloat64(nf sql.NullFloat64) float64 {
	if !nf.Valid {
		return 0
	}
	return nf.Float64
}

// Bool conversions
func SqlNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func SqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// NullBoolToPtr is already defined in conversions.go

func BoolToNullBool(b bool) sql.NullBool {
	return sql.NullBool{Bool: b, Valid: true}
}

func NullBoolToBool(nb sql.NullBool) bool {
	if !nb.Valid {
		return false
	}
	return nb.Bool
}

// Helper functions for common patterns

// StringPtr creates a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// UUIDPtr creates a pointer to a UUID
func UUIDPtr(u uuid.UUID) *uuid.UUID {
	return &u
}

// Int32Ptr creates a pointer to an int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Int64Ptr creates a pointer to an int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr creates a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr creates a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}
