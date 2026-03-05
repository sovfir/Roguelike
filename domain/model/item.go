package model

import (
	"rogue/infrastructure/constants"
)

type Item struct {
	ItemType  constants.EntityType
	Name      string
	Health    int
	MaxHealth int
	Agility   int
	Strength  int
	Value     int
}

func NewItem(
	itemType constants.EntityType,
	name string,
	health int,
	maxHealth int,
	agility int,
	strength int,
	value int,
) *Item {
	return &Item{
		ItemType:  itemType,
		Name:      name,
		Health:    health,
		MaxHealth: maxHealth,
		Agility:   agility,
		Strength:  strength,
		Value:     value,
	}
}
