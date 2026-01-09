package auth

import (
	"context"
	"slap-realtime/internal/store"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSqlRepository_GetUserByPhone(t *testing.T) {
	db := store.NewTestDB(t)
	defer db.Close()
	store.TruncateTables(t, db, "users")

	queries := store.New(db)
	repo := NewRepository(queries)

	// --- Seed data ---
	userID := uuid.New()
	phoneNumber := "+15551234567"
	_, err := db.Exec("INSERT INTO users (id, phone_number) VALUES ($1, $2)", userID, phoneNumber)
	assert.NoError(t, err)

	t.Run("found", func(t *testing.T) {
		foundID, err := repo.GetUserByPhone(context.Background(), phoneNumber)
		assert.NoError(t, err)
		assert.Equal(t, userID, foundID)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetUserByPhone(context.Background(), "+15550000000")
		assert.Error(t, err)
	})
}