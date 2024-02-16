package user

import (
	"context"

	"github.com/gocql/gocql"
)

type LoginInfoRepository struct {
	db *gocql.Session
}

func NewLoginInfoRepository(db *gocql.Session) *LoginInfoRepository {
	return &LoginInfoRepository{db: db}
}

func (r LoginInfoRepository) Insert(ctx context.Context, email, userId string) error {
	userID, err := gocql.ParseUUID(userId)
	if err != nil {
		return err
	}

	return r.db.Query("INSERT INTO login_info (email, user_id) VALUES (?, ?)", email, userID).WithContext(ctx).Exec()
}

func (r LoginInfoRepository) GetUserIDByEmail(ctx context.Context, email string) (string, error) {
	var userId string
	if err := r.db.Query("SELECT user_id FROM login_info WHERE email = ?", email).WithContext(ctx).Scan(&userId); err != nil {
		return "", err
	}

	return userId, nil
}

func (r LoginInfoRepository) IsUniqueEmail(ctx context.Context, email string) (bool, error) {
	var count int
	if err := r.db.Query("SELECT COUNT(*) FROM login_info WHERE email = ?", email).WithContext(ctx).Scan(&count); err != nil {
		return false, err
	}

	return count == 0, nil
}
