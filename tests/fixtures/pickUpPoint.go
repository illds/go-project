package fixtures

import (
	"GOHW-1/internal/model"
	"GOHW-1/tests/states"
)

type PickUpPointBuilder struct {
	instance *model.PickUpPoint
}

func PickUpPoint() *PickUpPointBuilder {
	return &PickUpPointBuilder{instance: &model.PickUpPoint{}}
}

func (b *PickUpPointBuilder) ID(v int64) *PickUpPointBuilder {
	b.instance.ID = v
	return b
}

func (b *PickUpPointBuilder) Name(v string) *PickUpPointBuilder {
	b.instance.Name = v
	return b
}

func (b *PickUpPointBuilder) Contact(v string) *PickUpPointBuilder {
	b.instance.Contact = v
	return b
}

func (b *PickUpPointBuilder) Address(v string) *PickUpPointBuilder {
	b.instance.Address = v
	return b
}

func (b *PickUpPointBuilder) P() *model.PickUpPoint {
	return b.instance
}

func (b *PickUpPointBuilder) V() model.PickUpPoint {
	return *b.instance
}

func (b *PickUpPointBuilder) Valid() *PickUpPointBuilder {
	return PickUpPoint().ID(states.PickUpPoint1ID).Name(states.PickUpPoint1Name).
		Address(states.PickUpPoint1Address).Contact(states.PickUpPoint1Contact)
}
