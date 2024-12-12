package types

import "github.com/google/uuid"

type Project struct {
	tableName   struct{}  `pg:"project"`
	ID          uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Tick        string    `json:"tick"`
	Type        string    `json:"type"`
	Image       string    `json:"image"`
	Information string    `json:"information"`
	Stage1      int64     `json:"stage_1"`
	Stage2      int64     `json:"stage_2"`
	Stage3      int64     `json:"stage_3"`
	Sort        uint64    `json:"sort"`
	CreatedAt   int64     `json:"created_at"`
	UpdateAt    int64     `json:"update_at"`
}
