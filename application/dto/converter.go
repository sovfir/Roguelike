package dto

import (
	"rogue/domain/model"
	"rogue/infrastructure/constants"
)

func GetGameInfo(gs *model.GameSession) DomainToViewDTO {
	gsCells := gs.Cells()
	gsItems := gs.Items()
	gsEnemies := gs.Enemies()
	cells := make([]CellInfoDTO, 0, len(gsCells))
	for key, val := range gsCells {
		if val.IsVisible() || val.IsVisited() {
			entity := constants.NONE
			if val.IsVisible() {
				if item, exists := gsItems[key]; exists {
					entity = item.ItemType
				}
				if enemy, exists := gsEnemies[key]; exists {
					entity = enemy.LooksLike()
				}
			}
			cellInfo := newCellInfoDTO(key.Row(), key.Col(), val.Ground(), entity)
			cells = append(cells, cellInfo)
		}
	}
	heroInfo := newHeroInfoDTO(
		gs.Hero.Coordinate.Row(),
		gs.Hero.Coordinate.Col(),
		gs.Hero.Health,
		gs.Hero.MaxHealth,
		gs.Hero.Agility,
		gs.Hero.Strength,
		gs.Hero.Gold,
		gs.Hero.IsSleeping,
	)

	status := constants.IN_GAME
	level := gs.LevelNumber
	if level > 21 {
		gs.SaveStats()
		status = constants.WIN

	}
	if gs.Hero.Health <= 0 {
		gs.SaveStats()
		status = constants.LOSE
	}

	weapon := ""
	if gs.Hero.Weapon != nil {
		weapon = gs.Hero.Weapon.Name
	}
	return newGameInfoDTO(cells, heroInfo, level, gs.Message, status, weapon)
}
