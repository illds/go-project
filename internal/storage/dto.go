package storage

import "time"

type OrderDTO struct {
	ID             int
	ClientID       int
	ExpirationDate time.Time
	Weight         float64
	Price          float64
	Packaging      string
	IsGiven        bool
	GivenTime      time.Time
}
