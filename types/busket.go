package types

import "github.com/google/uuid"

type Busket struct {
	ID       string
	Products []*BusketProduct
	UID      string
	Balance  float64
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
	ID        string
	UnitPrice float64
	UnitSize  int
}
