package types

import "github.com/google/uuid"

type BtcOrder struct {
	tableName struct{}  `pg:"btc_order"`
	ID        uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	Type      uint8
	TxId      string
	TxHash    string
	// Pending Inscribed NotPaid
	Status         string
	PayAddress     string
	ReceiveAddress string
	Files          []File `pg:"rel:has-many"`
	ClientId       string
	Count          uint64
	FeeRate        uint32
	MinerFee       uint64
	ServiceFee     uint64
	Amount         uint64
	//Data           OrderData `pg:"rel:has-many"`
	ExTime    int64
	CreatedAt int64
}

type Data struct {
	tableName      struct{}  `pg:"btc_data"`
	ID             uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	BtcOrderId     uuid.UUID
	OutputIndex    uint32 `json:"outputIndex"`
	PaymentAddress string `json:"paymentAddress"`
}

type File struct {
	tableName  struct{}  `pg:"btc_file"`
	ID         uuid.UUID `bun:"type:uuid,default:uuid_generate_v4()"`
	BtcOrderId uuid.UUID
	Index      uint64
	FileName   string
	DataUrl    string
	Size       uint64
	TxId       string
	Address    string
}
