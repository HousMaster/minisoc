package sqlite

import (
	"app/internal/storage"
	"app/internal/storage/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func (s *Storage) GetUser(username string) (*models.User, error) {
	const op = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT id,password FROM users WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	user := new(models.User)
	if err = stmt.QueryRow(username).Scan(&user.ID, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return user, nil
}

func (s *Storage) CreateUser(user *models.User) (int64, error) {
	const op = "storage.sqlite.CreateUser"

	smtp, err := s.db.Prepare("INSERT INTO users(username, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := smtp.Exec(user.Username, user.Password)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, storage.ErrUserExists
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserProfile(username string) (*models.UserProfile, error) {
	const op = "storage.sqlite.GetUserProfile"

	stmt, err := s.db.Prepare("SELECT description FROM users WHERE username = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	userProfile := new(models.UserProfile)
	userProfile.Username = username
	if err = stmt.QueryRow(username).Scan(&userProfile.Description); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if userProfile.Description == "" {
		return nil, storage.ErrUserDescriptionIsEmpty
	}

	return userProfile, nil
}

func (s *Storage) UpdateUserProfile(userProfile models.UserProfile) error {
	const op = "storage.sqlite.UpdateUserProfile"

	smtp, err := s.db.Prepare("UPDATE users SET description = ? WHERE username = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = smtp.Exec(userProfile.Description, userProfile.Username)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return storage.ErrUserExists
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
