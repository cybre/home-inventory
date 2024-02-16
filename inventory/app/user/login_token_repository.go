package user

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

type LoginTokenRepository struct {
	db *gocql.Session
}

func NewLoginTokenRepository(db *gocql.Session) *LoginTokenRepository {
	return &LoginTokenRepository{db: db}
}

func (r LoginTokenRepository) InsertToken(ctx context.Context, userId, token string) error {
	userID, err := gocql.ParseUUID(userId)
	if err != nil {
		return fmt.Errorf("failed to parse user id: %w", err)
	}

	tokenUUID, err := gocql.ParseUUID(token)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	return r.db.Query("INSERT INTO login_tokens (login_token, user_id) VALUES (?, ?)", tokenUUID, userID).WithContext(ctx).Exec()
}

func (r LoginTokenRepository) VerifyToken(ctx context.Context, token string) (string, error) {
	tokenUUID, err := gocql.ParseUUID(token)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	var userId string
	if err := r.db.Query("SELECT user_id FROM login_tokens WHERE login_token = ?", tokenUUID).WithContext(ctx).Scan(&userId); err != nil {
		return "", fmt.Errorf("failed to get user id by token: %w", err)
	}

	return userId, nil
}
