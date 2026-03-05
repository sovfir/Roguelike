package usecases

import (
	"math/rand/v2"

	"rogue/domain/model"
	"rogue/domain/service"
	"rogue/infrastructure/constants"
)

func EnemiesMove(gs *model.GameSession) {
	enemies := gs.Enemies()
	sliceEnemies := make([]model.Coordinate, 0, len(enemies))
	for key := range enemies {
		sliceEnemies = append(sliceEnemies, key)
	}

	for _, coordinate := range sliceEnemies {
		enemy := enemies[coordinate]
		if !enemy.InPursuing() {
			if enemy.InAgrDistance(coordinate, gs.Hero.Coordinate) {
				switch enemy.Type() {
				case constants.VAMPIRE:
					vampire := enemy.(model.Vampire)
					vampire.StartPursuing()
					enemy = vampire
				case constants.ZOMBIE:
					zombie := enemy.(model.Zombie)
					zombie.StartPursuing()
					enemy = zombie
				case constants.OGRE:
					ogre := enemy.(model.Ogre)
					ogre.StartPursuing()
					enemy = ogre
				case constants.GHOST:
					ghost := enemy.(model.Ghost)
					ghost.StartPursuing()
					enemy = ghost
				case constants.SNAKE:
					snake := enemy.(model.Snake)
					snake.StartPursuing()
					enemy = snake
				case constants.MIMIC:
					mimic := enemy.(model.Mimic)
					mimic.StartPursuing()
					enemy = mimic
				}
			}
		}

		newCoordinate := coordinate
		isRegularMove := true
		if enemy.InPursuing() {
			newCoordinate = enemy.PathFinder(coordinate, gs.Hero.Coordinate, gs)
			if newCoordinate == coordinate {
				isRegularMove = true
			} else {
				isRegularMove = false
			}
			if newCoordinate == gs.Hero.Coordinate {
				newCoordinate = coordinate
				if enemy.Type() == constants.OGRE {
					ogre := enemy.(model.Ogre)
					if ogre.Resting {
						ogre.Resting = false
					} else {
						ogre.Attack(gs)
						ogre.Resting = true
					}
					enemy = ogre
				} else {
					enemy.Attack(gs)
				}
			}
		}
		if isRegularMove {
			room, _ := service.RoomByCoordinate(coordinate, gs)
			newCoordinate = enemy.RegularMove(room, coordinate, gs)
			if enemy.Type() == constants.GHOST {
				ghost := enemy.(model.Ghost)
				ghost.Visibility = (rand.IntN(2) == 1)
				enemy = ghost
			}
		}
		gs.MoveEnemy(coordinate, newCoordinate, enemy)
	}
}
