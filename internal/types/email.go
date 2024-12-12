package types

import (
	"github.com/google/uuid"
	"time"
)

type Email struct {
	tableName struct{}  `pg:"btc_email"`
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
