package service

import (
	"math/rand/v2"
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

func UseItem(item *model.Item, gs *model.GameSession) {
	switch item.ItemType {
	case constants.WEAPON:
		if gs.Hero.Weapon != nil {
			// выбрасывание текущего предмета
			neighbours := make([]model.Coordinate, 0, 8)
			center := gs.Hero.Coordinate
			for i := -1; i <= 1; i++ {
				for j := -1; j <= 1; j++ {
					if i == 0 && j == 0 {
						continue
					}
					neighbour := model.NewCoordinate(center.Row()+i, center.Col()+j)

					if cell, exists := gs.Cells()[neighbour]; exists && !occupied(neighbour, cell, gs) {
						neighbours = append(neighbours, neighbour)
					}
				}
			}
			place := neighbours[rand.IntN(len(neighbours))]
			gs.AddItem(place, gs.Hero.Weapon)
		}

		gs.Hero.Weapon = item
	case constants.ELIXIR:
		gs.Hero.IncreaseStats(item)
		gs.SessionStats.ElixirsConsumed++
		time := constants.EFFECT_BASE_TIME - gs.LevelNumber/2
		newEffect := model.NewEffect(item, time)
		gs.Effects[len(gs.Effects)] = newEffect
	case constants.SCROLL:
		gs.Hero.IncreaseStats(item)
		gs.SessionStats.ScrollsConsumed++
	default:
		gs.Hero.IncreaseStats(item)
		gs.SessionStats.FoodConsumed++
	}
}

func occupied(coordinate model.Coordinate, cell model.Cell, gs *model.GameSession) bool {
	if _, exists := gs.Items()[coordinate]; exists {
		return true
	}
	if cell.Ground() == constants.WALL || cell.Ground() == constants.EXIT {
		return true
	}
	// добавить проверку на монстра
	return false
}
