package storage

import (
	"app/internal/storage/models"
)

type Message interface {
	GetMessages() ([]models.Message, error)
	CreateMessage(*models.Message) (int64, error)
}
