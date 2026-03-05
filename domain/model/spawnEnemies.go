package model

import (
	"math/rand/v2"

	"rogue/infrastructure/constants"
)

func SpawnEnemies(levelNumber int, level *Level, balancer float64) {
	enemyTypes := []constants.EntityType{
		constants.VAMPIRE,
		constants.ZOMBIE,
		constants.OGRE,
		constants.GHOST,
		constants.SNAKE,
		constants.MIMIC,
	}
	itemTypes := []constants.EntityType{
		constants.FOOD,
		constants.ELIXIR,
		constants.SCROLL,
		constants.WEAPON,
	}
	countBeforeBalance := float64(constants.BASE_COUNT_ENEMY_IN_DUNGEON + levelNumber/constants.LEVEL_UPDATE_DIFFICULTY)
	maxCount := int(countBeforeBalance * balancer)
	countSpawned := 0

	for countSpawned < maxCount {
		roomIndex := rand.IntN(len(level.rooms))
		if roomIndex == level.StartRoomIndex {
			continue
		}
		coordinate := level.rooms[roomIndex].GetRandomCoordinate()
		if _, exists := level.items[coordinate]; exists {
			continue
		}
		if _, exists := level.enemies[coordinate]; exists {
			continue
		}
		enemyType := enemyTypes[rand.IntN(len(enemyTypes))]
		var enemy Enemy
		levelDifficulty := (2 * levelNumber) / constants.LEVEL_UPDATE_DIFFICULTY
		baseHealth := int(float64(constants.BASE_ENEMY_HEALTH+levelDifficulty) * balancer)
		baseAgility := int(float64(constants.BASE_ENEMY_AGILITY+levelDifficulty) * balancer)
		baseStrength := int(float64(constants.BASE_ENEMY_STRENGTH+levelDifficulty) * balancer)
		varHealth := constants.VARIABLE_ENEMY_HEALTH + levelDifficulty
		varAgility := constants.VARIABLE_ENEMY_AGILITY + levelDifficulty
		varStrength := constants.VARIABLE_ENEMY_STRENGTH + levelDifficulty
		switch enemyType {
		case constants.VAMPIRE:
			enemy = Vampire{
				Monster: Monster{
					MonsterType: constants.VAMPIRE,
					Name:        "Vampire",
					Health:      baseHealth + constants.HIGH*varHealth,
					Agility:     baseAgility + constants.HIGH*varAgility,
					Strength:    baseStrength + constants.MEDIUM*varStrength,
					Hostility:   constants.HIGH,
					InBattle:    false,
				},
				FirstHit: true,
			}
		case constants.ZOMBIE:
			enemy = Zombie{
				Monster: Monster{
					MonsterType: constants.ZOMBIE,
					Name:        "Zombie",
					Health:      baseHealth + constants.HIGH*varHealth,
					Agility:     baseAgility + constants.LOW*varAgility,
					Strength:    baseStrength + constants.MEDIUM*varStrength,
					Hostility:   constants.MEDIUM,
					InBattle:    false,
				},
			}
		case constants.OGRE:
			enemy = Ogre{
				Monster: Monster{
					MonsterType: constants.OGRE,
					Name:        "Ogre",
					Health:      baseHealth + constants.VERY_HIGH*varHealth,
					Agility:     baseAgility + constants.LOW*varAgility,
					Strength:    baseStrength + constants.VERY_HIGH*varStrength,
					Hostility:   constants.MEDIUM,
					InBattle:    false,
				},
				Resting: false,
			}
		case constants.GHOST:
			enemy = Ghost{
				Monster: Monster{
					MonsterType: constants.GHOST,
					Name:        "Ghost",
					Health:      baseHealth + constants.LOW*varHealth,
					Agility:     baseAgility + constants.HIGH*varAgility,
					Strength:    baseStrength + constants.LOW*varStrength,
					Hostility:   constants.LOW,
					InBattle:    false,
				},
				Visibility: true,
			}
		case constants.SNAKE:
			enemy = Snake{
				Monster: Monster{
					MonsterType: constants.SNAKE,
					Name:        "Snake-Mage",
					Health:      baseHealth + constants.LOW*varHealth,
					Agility:     baseAgility + constants.VERY_HIGH*varAgility,
					Strength:    baseStrength + constants.LOW*varStrength,
					Hostility:   constants.HIGH,
					InBattle:    false,
				},
			}
		case constants.MIMIC:
			enemy = Mimic{
				Monster: Monster{
					MonsterType: constants.MIMIC,
					Name:        "Mimic",
					Health:      baseHealth + constants.HIGH*varHealth,
					Agility:     baseAgility + constants.HIGH*varAgility,
					Strength:    baseStrength + constants.LOW*varStrength,
					Hostility:   constants.LOW,
					InBattle:    false,
				},
				ItemMimic: itemTypes[rand.IntN(len(itemTypes))],
			}
		}
		level.enemies[coordinate] = enemy
		countSpawned++
	}
}
