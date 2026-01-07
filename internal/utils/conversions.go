package utils

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// InterfaceSliceToStringSlice converts a slice of interface{} to a slice of string
// This is commonly needed when working with sqlc-generated code for PostgreSQL enums
func InterfaceSliceToStringSlice(slice []interface{}) []string {
	result := make([]string, len(slice))
	for i, v := range slice {
		if str, ok := v.(string); ok {
			result[i] = str
		} else {
			// Handle the case where the conversion fails
			result[i] = ""
		}
	}
	return result
}

// InterfaceToString converts an interface{} to a string
// This is commonly needed when working with sqlc-generated code for PostgreSQL enums
func InterfaceToString(i interface{}) string {
	if str, ok := i.(string); ok {
		return str
	}
	return ""
}

// SafeInterfaceToString converts an interface{} to a string with safety check
func SafeInterfaceToString(i interface{}) string {
	if i == nil {
		return ""
	}
	if str, ok := i.(string); ok {
		return str
	}
	return ""
}

// SafeInterfaceSliceToStringSlice converts a slice of interface{} to a slice of string with safety checks
func SafeInterfaceSliceToStringSlice(slice []interface{}) []string {
	if slice == nil {
		return []string{}
	}
	result := make([]string, len(slice))
	for i, v := range slice {
		if v == nil {
			result[i] = ""
		} else if str, ok := v.(string); ok {
			result[i] = str
		} else {
			result[i] = ""
		}
	}
	return result
}

// NullUUIDToPtr converts a uuid.NullUUID to a *uuid.UUID
func NullUUIDToPtr(n uuid.NullUUID) *uuid.UUID {
	if n.Valid {
		return &n.UUID
	}
	return nil
}

// PtrToNullUUID converts a *uuid.UUID to a uuid.NullUUID
func PtrToNullUUID(u *uuid.UUID) uuid.NullUUID {
	if u == nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: *u, Valid: true}
}

// NullBoolToPtr converts a sql.NullBool to a *bool
func NullBoolToPtr(n sql.NullBool) *bool {
	if n.Valid {
		return &n.Bool
	}
	return nil
}

// NullFloat64ToFloat32Ptr converts a sql.NullFloat64 to a *float32
func NullFloat64ToFloat32Ptr(n sql.NullFloat64) *float32 {
	if n.Valid {
		f := float32(n.Float64)
		return &f
	}
	return nil
}

// NullTimeToPtr converts a sql.NullTime to a *time.Time
func NullTimeToPtr(n sql.NullTime) *time.Time {
	if n.Valid {
		return &n.Time
	}
	return nil
}

// StringSliceToInterfaceSlice converts a slice of string to a slice of interface{}
// This is commonly needed when working with sqlc-generated code for PostgreSQL enums
func StringSliceToInterfaceSlice(slice []string) []interface{} {
	result := make([]interface{}, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}
