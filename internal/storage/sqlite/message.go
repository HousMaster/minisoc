package sqlite

import (
	"app/internal/storage/models"
	"fmt"
)

func (s *Storage) messageInit() error {
	const op = "storage.sqlite.messageInit"

	stmt, err := s.db.Prepare(`
	CREATE TABLE IF NOT EXISTS messages(
		id INTEGER PRIMARY KEY,
		from_id BIGINT NOT NULL,
		text TEXT,
		create_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP);
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetMessages() ([]models.Message, error) {

	const op = "storage.sqlite.GetMessages"

	stmt, err := s.db.Prepare("SELECT * FROM messages")
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: execute query: %w", op, err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {

		var id int64
		var fromID int64
		var text string
		var mCreateTime string

		err := rows.Scan(&id, &fromID, &text, &mCreateTime)

		if err != nil {
			return nil, fmt.Errorf("%s: scan row: %w", op, err)
		}

		messages = append(messages, models.Message{
			ID:     id,
			FromID: fromID,
			Text:   text,
		})

	}
	return messages, nil
}

func (s *Storage) CreateMessage(m *models.Message) (int64, error) {

	const op = "storage.sqlite.CreateMessage"

	smtp, err := s.db.Prepare("INSERT INTO messages(from_id, text) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := smtp.Exec(m.FromID, m.Text)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
