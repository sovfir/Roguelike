package model

import (
	"math/rand/v2"
	"strconv"

	"rogue/infrastructure/constants"
)

type Hero struct {
	Coordinate   Coordinate
	Health       int
	MaxHealth    int
	Agility      int
	Strength     int
	Weapon       *Item
	HasFirstKey  bool
	HasSecondKey bool
	IsSleeping   bool
	Gold         int
}

func NewHero(coordinate Coordinate, health, maxHealth, agility, strength int) Hero {
	return Hero{
		Coordinate: coordinate,
		Health:     health,
		MaxHealth:  maxHealth,
		Agility:    agility,
		Strength:   strength,
		Weapon:     nil,
	}
}

func (h *Hero) IncreaseStats(item *Item) {
	h.Health += item.Health
	h.MaxHealth += item.MaxHealth
	h.Agility += item.Agility
	h.Strength += item.Strength

	if h.Health > h.MaxHealth {
		h.Health = h.MaxHealth
	}
}

func (h *Hero) DecreaseStats(item *Item) {
	h.Health -= item.Health
	h.MaxHealth -= item.MaxHealth
	h.Agility -= item.Agility
	h.Strength -= item.Strength

	if h.Health <= 0 {
		h.Health = 1
	}
}

func (h Hero) Attack(gs *GameSession, coordinate Coordinate) {
	enemy := gs.level.enemies[coordinate]

	if enemy.Type() == constants.VAMPIRE {
		vampire := enemy.(Vampire)
		if vampire.FirstHit {
			vampire.FirstHit = false
			gs.level.enemies[coordinate] = vampire
			gs.Message += " Hero missed"
			return
		}
	}

	if !enemy.HitCheck(h.Agility) {
		gs.Message += " Hero missed"
		return
	}

	damageCoefficient := 1 - constants.DAMAGE_DISPERSION/2 + rand.Float64()*constants.DAMAGE_DISPERSION
	weaponStrength := 0
	if h.Weapon != nil {
		weaponStrength = h.Weapon.Strength
	}
	heroDamage := float64(h.Strength+weaponStrength) * damageCoefficient

	critChance := constants.BASE_CRIT_CHANCE * float64(h.Agility) / float64(enemy.GetAgility())
	if rand.Float64() < critChance {
		heroDamage *= constants.CRIT_DAMAGE_COEFFICIENT
	}

	gs.Message += " Hero dealt " + strconv.Itoa(int(heroDamage)) + " damage"

	var value int

	switch enemy.Type() {
	case constants.VAMPIRE:
		vampire := enemy.(Vampire)
		value = vampire.TakeDamage(int(heroDamage))
		enemy = vampire
	case constants.ZOMBIE:
		zombie := enemy.(Zombie)
		value = zombie.TakeDamage(int(heroDamage))
		enemy = zombie
	case constants.OGRE:
		ogre := enemy.(Ogre)
		value = ogre.TakeDamage(int(heroDamage))
		enemy = ogre
	case constants.GHOST:
		ghost := enemy.(Ghost)
		value = ghost.TakeDamage(int(heroDamage))
		enemy = ghost
	case constants.SNAKE:
		snake := enemy.(Snake)
		value = snake.TakeDamage(int(heroDamage))
		enemy = snake
	case constants.MIMIC:
		mimic := enemy.(Mimic)
		value = mimic.TakeDamage(int(heroDamage))
		mimic.StartPursuing()
		enemy = mimic
	}

	if value != 0 {
		gs.Message += " " + enemy.GetName() + " is dead"
		// Сохраняем статистику
		gs.SessionStats.EnemysKilled++
		delete(gs.level.enemies, coordinate)
		gs.level.items[coordinate] = &Item{
			ItemType: constants.TREASURE,
			Name:     strconv.Itoa(value) + " gold",
			Value:    value,
		}
	} else {
		gs.level.enemies[coordinate] = enemy
	}
}
