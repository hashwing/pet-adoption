package common

import (
	"github.com/gofrs/uuid"
)

func NewUUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
