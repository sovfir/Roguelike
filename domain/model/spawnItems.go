package model

import (
	"math/rand/v2"
	"strconv"

	"rogue/infrastructure/constants"
)

type stats int

const (
	MAX_HEALTH stats = iota
	AGILITY
	STRENGTH
)

func SpawnItems(hero Hero, levelNumber int, level *Level, balancer float64) {
	countBeforeBalance := float64(constants.BASE_COUNT_ITEM_IN_DUNGEON - levelNumber/constants.LEVEL_UPDATE_DIFFICULTY)
	maxCount := int(countBeforeBalance / balancer)
	countSpawned := 0
	itemTypes := []constants.EntityType{
		constants.FOOD,
		constants.ELIXIR,
		constants.SCROLL,
		constants.WEAPON,
	}
	for countSpawned < maxCount {
		roomIndex := rand.IntN(len(level.rooms))
		if roomIndex == level.StartRoomIndex {
			continue
		}
		coordinate := level.rooms[roomIndex].GetRandomCoordinate()
		if cell := level.cells[coordinate]; cell.ground == constants.EXIT {
			continue
		}
		if _, exists := level.items[coordinate]; exists {
			continue
		}
		itemType := itemTypes[rand.IntN(len(itemTypes))]
		var item *Item
		switch itemType {
		case constants.FOOD:
			item = createFood(hero)
		case constants.ELIXIR:
			item = createElixir(hero)
		case constants.SCROLL:
			item = createScroll(hero)
		case constants.WEAPON:
			item = createWeapon(levelNumber)
		}
		level.items[coordinate] = item
		countSpawned++
	}
}

func createFood(hero Hero) *Item {
	names := constants.FOOD_NAMES

	valuePercentage := randomInt(constants.MIN_FOOD_PERCENT, constants.MAX_FOOD_PERCENT)
	health := hero.MaxHealth * valuePercentage / 100
	name := names[rand.IntN(len(names))] + " +" + strconv.Itoa(health)
	return NewItem(constants.FOOD, name, health, 0, 0, 0, 0)
}

func createElixir(hero Hero) *Item {
	names := constants.ELIXIR_NAMES
	stats := []stats{MAX_HEALTH, AGILITY, STRENGTH}
	stat := stats[rand.IntN(len(stats))]

	valuePercentage := randomInt(constants.MIN_ELIXIR_PERCENT, constants.MAX_ELIXIR_PERCENT)

	health := 0
	agility := 0
	strength := 0
	name := ""
	switch stat {
	case MAX_HEALTH:
		health = hero.MaxHealth * valuePercentage / 100
		name = names[rand.IntN(len(names))] + " +" + strconv.Itoa(health) + " health"
	case AGILITY:
		agility = hero.Agility * valuePercentage / 100
		name = names[rand.IntN(len(names))] + " +" + strconv.Itoa(agility) + " agility"
	case STRENGTH:
		strength = hero.Strength * valuePercentage / 100
		name = names[rand.IntN(len(names))] + " +" + strconv.Itoa(strength) + " strength"
	}
	return NewItem(constants.ELIXIR, name, health, health, agility, strength, 0)
}

func createScroll(hero Hero) *Item {
	names := constants.SCROLL_NAMES
	stats := []stats{MAX_HEALTH, AGILITY, STRENGTH}
	stat := stats[rand.IntN(len(stats))]

	valuePercentage := randomInt(constants.MIN_SCROLL_PERCENT, constants.MAX_SCROLL_PERCENT)
	health := 0
	agility := 0
	strength := 0
	name := ""
	switch stat {
	case MAX_HEALTH:
		health = hero.MaxHealth * valuePercentage / 100
		name = names[rand.IntN(len(names))] + " +" + strconv.Itoa(health) + " maximum health"
	case AGILITY:
		agility = hero.Agility * valuePercentage / 100
		name = names[rand.IntN(len(names))] + " +" + strconv.Itoa(agility) + " agility"
	case STRENGTH:
		strength = hero.Strength * valuePercentage / 100
		name = names[rand.IntN(len(names))] + " +" + strconv.Itoa(strength) + " strength"
	}
	return NewItem(constants.SCROLL, name, health, health, agility, strength, 0)
}

func createWeapon(levelNumber int) *Item {
	names := constants.WEAPON_NAMES

	strength := randomInt(constants.MIN_WEAPON_STRENGTH, constants.MAX_WEAPON_STRENGTH+levelNumber/2)
	name := names[rand.IntN(len(names))] + " +" + strconv.Itoa(strength)
	return NewItem(constants.WEAPON, name, 0, 0, 0, strength, 0)
}

func randomInt(min, max int) int {
	return min + rand.IntN(max-min+1)
}
