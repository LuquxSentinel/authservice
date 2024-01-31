package types

import "github.com/google/uuid"

type Busket struct {
	ID       string           `json:"id" bson:"id"`
	Products []*BusketProduct `json:"products" bson:"products"`
	UID      string           `json:"uid" bson:"uid"`
	Balance  float64          `json:"balace" bson:"balance"`
}

func NewBusket(uid string) *Busket {
	return &Busket{
		ID:       uuid.New().String(),
		Products: make([]*BusketProduct, 0),
		UID:      uid,
		Balance:  0.00,
	}
}

type BusketProduct struct {
	ID        string  `json:"id" bson:"id"`
	UnitPrice float64 `json:"unit_price" bson:"unit_price"`
	UnitSize  int     `json:"unit_size" bson:"unit_size"`
}
