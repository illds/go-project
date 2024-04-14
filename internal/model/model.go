package model

import (
	"errors"
	"time"
)

var ErrObjectNotFound = errors.New("not found")

type Order struct {
	ID             int
	ClientID       int
	ExpirationDate time.Time
	Weight         float64 // In kg
	Price          float64
	Packaging      string
}

type PickUpPoint struct {
	ID      int64  `db:"id"`
	Name    string `db:"name"`
	Address string `db:"address"`
	Contact string `db:"contact"`
}

const (
	Package string = "package"
	Carton  string = "carton"
	Film    string = "film"
)

type PackagingRule struct {
	MaxWeight float64 // Max weight in kg, 0 - no limit
	ExtraCost float64 // Extra cost in rubles
}

func GetPackagingRules() map[string]PackagingRule {
	return map[string]PackagingRule{
		Package: {MaxWeight: 10, ExtraCost: 5},
		Carton:  {MaxWeight: 30, ExtraCost: 20},
		Film:    {ExtraCost: 1},
	}
}
