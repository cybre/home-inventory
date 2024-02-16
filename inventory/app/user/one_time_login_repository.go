package user

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
)

type OneTimeLoginRepository struct {
	db *gocql.Session
}

func NewOneTimeLoginRepository(db *gocql.Session) *OneTimeLoginRepository {
	return &OneTimeLoginRepository{db: db}
}

func (r OneTimeLoginRepository) InsertToken(ctx context.Context, userId, token string) error {
	userID, err := gocql.ParseUUID(userId)
	if err != nil {
		return fmt.Errorf("failed to parse user id: %w", err)
	}

	tokenUUID, err := gocql.ParseUUID(token)
	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	return r.db.Query("INSERT INTO login_tokens (login_token, user_id) VALUES (?, ?)", userID, tokenUUID).WithContext(ctx).Exec()
}
